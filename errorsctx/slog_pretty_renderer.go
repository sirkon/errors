package errorsctx

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/sirkon/errors"
)

// NodeKind определяет тип конечного значения для точечной раскраски
type NodeKind int

const (
	KindGroup NodeKind = iota
	KindArray
	KindString
	KindInlineArray // Для компактного вывода []byte в одну строку
	KindNumber
	KindBool
	KindNull
	KindLocation
	KindErrorText
	KindErrorNode
	KindStackTrace
)

// TreeNode — единый элемент нашего сквозного промежуточного дерева
type TreeNode struct {
	Key        string
	Value      string
	Kind       NodeKind
	Children   []*TreeNode
	RawDisplay bool
	IsHex      bool
}

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

type SlogPrettyRenderer struct {
	opts     slog.HandlerOptions
	preAttrs []slog.Attr
	dst      io.Writer
	color    *prettyWriterColorProfile
	hexLimit int
}

// NewSlogPrettyRenderer creates pretty output slog.Handler.
//
//   - isDark applies respective color profile
//   - hexLimit truncates binary data longer than the limit value. Here -1 disables this functionality and 0 is
//     interpreted as 32.
func NewSlogPrettyRenderer(dst io.Writer, opts *slog.HandlerOptions, isDark bool, hexLimit int) *SlogPrettyRenderer {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	var profile *prettyWriterColorProfile
	if isDark {
		profile = newPrettyWriterColorProfileDark()
	} else {
		profile = newPrettyWriterColorProfileLight()
	}
	if hexLimit == 0 {
		hexLimit = 32
	}
	return &SlogPrettyRenderer{
		opts:     *opts,
		dst:      dst,
		color:    profile,
		hexLimit: hexLimit,
	}
}

func (h *SlogPrettyRenderer) Enabled(_ context.Context, level slog.Level) bool {
	if h.opts.Level != nil {
		return h.opts.Level.Level() <= level
	}
	return slog.LevelInfo <= level
}

func (h *SlogPrettyRenderer) Handle(_ context.Context, r slog.Record) error {
	bufPtr := bufPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]

	defer func() {
		buf = append(buf, '\n')
		_, _ = h.dst.Write(buf)
		*bufPtr = buf
		bufPool.Put(bufPtr)
	}()

	// 1. Отрисовка Времени
	buf = append(buf, h.color.time...)
	buf = r.Time.AppendFormat(buf, "2006-01-02 15:04:05.000")
	buf = append(buf, h.color.reset...)
	buf = append(buf, ' ')

	// 2. Отрисовка Уровня лога
	switch r.Level {
	case slog.LevelDebug:
		buf = append(buf, h.color.debug+"DEBUG"...)
	case slog.LevelInfo:
		buf = append(buf, h.color.info+"INFO"...)
	case slog.LevelWarn:
		buf = append(buf, h.color.warn+"WARN"...)
	case slog.LevelError:
		buf = append(buf, h.color.error+"ERROR"...)
	default:
		buf = append(buf, r.Level.String()...)
	}
	buf = append(buf, h.color.reset...)
	buf = append(buf, ' ')

	// 3. Сообщение лога
	buf = append(buf, h.color.bold...)
	buf = append(buf, r.Message...)
	buf = append(buf, h.color.reset...)

	forceTree := false
	hasMultilineString := false
	hasInternalJSON := false // <-- Добавить флаг

	rawAttrs := make([]slog.Attr, 0, len(h.preAttrs)+r.NumAttrs())

	processRaw := func(a slog.Attr) {
		val := a.Value.Resolve()
		if val.Kind() == slog.KindGroup {
			g := val.Group()
			if len(g) == 1 && g[0].Key == "__slog_force_tree__" {
				forceTree = true
				return
			}
		}
		if val.Kind() == slog.KindString {
			valStr := val.String()
			if strings.Contains(valStr, "\n") {
				hasMultilineString = true
			}
			// Эвристика: Если внутри плоской строки прилетел JSON, требуем дерево
			if len(valStr) > 1 && (valStr[0] == '{' || valStr[0] == '[') {
				hasInternalJSON = true
			}
		}
		rawAttrs = append(rawAttrs, a)
	}

	for _, a := range h.preAttrs {
		processRaw(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		processRaw(a)
		return true
	})

	// Сценарий 1: Мало контекста -> Компактный однострочный JSON
	if !forceTree && !hasMultilineString && !hasInternalJSON && len(rawAttrs) <= 3 {
		hasGroupsOrComplex := false
		for _, a := range rawAttrs {
			k := a.Value.Resolve().Kind()
			if k == slog.KindGroup || k == slog.KindAny {
				hasGroupsOrComplex = true
				break
			}
		}

		if !hasGroupsOrComplex {
			if len(rawAttrs) == 0 {
				return nil
			}
			buf = append(buf, ' ')
			buf = append(buf, h.color.stdots...)
			buf = append(buf, '{')
			buf = append(buf, h.color.reset...)

			for i, a := range rawAttrs {
				if i > 0 {
					buf = append(buf, h.color.stdots...)
					buf = append(buf, ", "...)
					buf = append(buf, h.color.reset...)
				}
				buf = append(buf, h.color.key...)
				buf = append(buf, '"')
				buf = append(buf, a.Key...)
				buf = append(buf, '"')
				buf = append(buf, h.color.reset...)
				buf = append(buf, h.color.stdots...)
				buf = append(buf, ": "...)
				buf = append(buf, h.color.reset...)

				buf = append(buf, h.color.ctx...)
				val := a.Value.Resolve()
				if val.Kind() == slog.KindString {
					buf = append(buf, '"')
					buf = appendRawSlogValue(buf, val)
					buf = append(buf, '"')
				} else {
					buf = appendRawSlogValue(buf, val)
				}
				buf = append(buf, h.color.reset...)
			}
			buf = append(buf, h.color.stdots...)
			buf = append(buf, '}')
			buf = append(buf, h.color.reset...)
			return nil
		}
	}

	// Сценарий 2: Построение Единого Промежуточного Дерева (IR) для всего контекста
	rootNodes := make([]*TreeNode, 0, len(rawAttrs))
	for _, a := range rawAttrs {
		rootNodes = append(rootNodes, h.buildIRTree(a.Key, a.Value))
	}

	// Линейный рендеринг готового IR-графа
	buf = h.renderIRTree(buf, rootNodes, []bool{}, false)
	return nil
}

// buildIRTree — Фабрика сквозного IR-дерева с интегрированными эвристиками
func (h *SlogPrettyRenderer) buildIRTree(key string, val slog.Value) *TreeNode {
	resolved := val.Resolve()
	node := &TreeNode{Key: key}

	// 1. Группы slog.Group
	if resolved.Kind() == slog.KindGroup {
		node.Kind = KindGroup
		for _, subAttr := range resolved.Group() {
			node.Children = append(node.Children, h.buildIRTree(subAttr.Key, subAttr.Value))
		}
		return node
	}

	// 2. ЧЕСТНЫЙ ПЕРЕХВАТ ОШИБКИ БЕЗ ЭВРИСТИК
	if resolved.Kind() == slog.KindAny {
		if e, ok := resolved.Any().(error); ok {
			if node.Key == "" || node.Key == "!BADKEY" {
				node.Key = "err"
			}

			// Сам корневой узел ошибки размечаем как KindErrorNode
			node.Kind = KindErrorNode

			var err *errors.Error
			if er, ok := e.(*errors.Error); ok {
				err = er
			} else {
				err, _ = errors.AsType[*errors.Error](e)
			}

			// Текст ошибки горит красным
			node.Children = append(node.Children, &TreeNode{Key: "@text", Value: e.Error(), Kind: KindErrorText})

			// Блок контекста ошибки
			ctxNode := &TreeNode{Key: "@context", Kind: KindErrorNode}

			if err != nil {
				// Если это наша родная ошибка sirkon/errors, раскладываем дерево
				for _, subAttr := range errors.SLogTreeContext(err) {
					ctxNode.Children = append(ctxNode.Children, h.buildIRTree(subAttr.Key, subAttr.Value))
				}
			} else {
				// Если это чужая ошибка (foreign error), контекст пустой, но узел создаем
				ctxNode.Value = "{}"
				ctxNode.Kind = KindNull
			}

			node.Children = append(node.Children, ctxNode)
			return node
		}
	}

	// Извлекаем строковое значение для текстовых эвристик (стектрейсы, локации, JSON)
	var rawStr string
	switch resolved.Kind() {
	case slog.KindString:
		rawStr = resolved.String()
	case slog.KindAny:
		if s, ok := resolved.Any().(string); ok {
			rawStr = s
		} else if b, ok := resolved.Any().([]byte); ok {
			rawStr = string(b)
		}
	}

	// 2. Железобетонная эвристика: Распознаем стектрейс Go по структуре текста
	if (key == "stacktrace" || key == "stack" || strings.Contains(rawStr, "goroutine ")) && strings.Contains(rawStr, "\n") {
		node.Kind = KindStackTrace
		node.Value = rawStr
		return node
	}

	// 3. Эвристика: Перехват ошибок пакета sirkon/errors
	if resolved.Kind() == slog.KindAny && node.Kind != KindStackTrace {
		if e, ok := resolved.Any().(error); ok {
			var err *errors.Error
			if er, ok := e.(*errors.Error); ok {
				err = er
			} else {
				err, _ = errors.AsType[*errors.Error](e)
			}

			if err != nil {
				if node.Key == "" || node.Key == "!BADKEY" {
					node.Key = "err"
				}
				node.Kind = KindGroup
				node.Children = append(node.Children, &TreeNode{Key: "@text", Value: e.Error(), Kind: KindErrorText})

				ctxNode := &TreeNode{Key: "@context", Kind: KindGroup}
				for _, subAttr := range errors.SLogTreeContext(err) {
					ctxNode.Children = append(ctxNode.Children, h.buildIRTree(subAttr.Key, subAttr.Value))
				}
				node.Children = append(node.Children, ctxNode)
				return node
			}
		}
	}

	// 4. Эвристика: Парсинг вложенных JSON / Map структур
	var anyObj any
	isJSON := false

	var jsonData []byte
	if len(rawStr) > 1 && (rawStr[0] == '{' || rawStr[0] == '[') {
		if json.Unmarshal([]byte(rawStr), &anyObj) == nil {
			isJSON = true
			jsonData = []byte(rawStr)
		}
	} else if resolved.Kind() == slog.KindAny {
		if jsonBytes, err := json.Marshal(resolved.Any()); err == nil && len(jsonBytes) > 1 {
			if jsonBytes[0] == '{' || jsonBytes[0] == '[' {
				if json.Unmarshal(jsonBytes, &anyObj) == nil {
					isJSON = true
				}
			}
			jsonData = jsonBytes
		}
	}

	if isJSON {
		return h.buildIRTreeFromJSONBytes(key, jsonData)
	}

	// 5. Дефолтная обработка примитивов верхнего уровня slog
	if key == "@location" || strings.Contains(rawStr, ".go:") {
		node.Kind = KindLocation
		node.Value = rawStr
		return node
	}

	switch resolved.Kind() {
	case slog.KindBool:
		node.Kind = KindBool
	case slog.KindInt64, slog.KindUint64, slog.KindFloat64, slog.KindDuration:
		node.Kind = KindNumber
	default:
		node.Kind = KindString
	}

	buf := make([]byte, 0, 64)
	buf = appendRawSlogValue(buf, resolved)
	node.Value = string(buf)
	return node
}

func (h *SlogPrettyRenderer) buildIRTreeFromJSONBytes(key string, data []byte) *TreeNode {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	// Передаем управление рекурсивному токен-парсеру
	return h.parseJSONToken(key, dec)
}

func (h *SlogPrettyRenderer) parseJSONToken(key string, dec *json.Decoder) *TreeNode {
	t, err := dec.Token()
	if err != nil {
		return &TreeNode{Key: key, Kind: KindNull, Value: "null"}
	}

	delim, ok := t.(json.Delim)
	if !ok {
		// Это примитив (string, number, bool, null)
		node := &TreeNode{Key: key}
		switch v := t.(type) {
		case bool:
			node.Kind = KindBool
			node.Value = strconv.FormatBool(v)
		case json.Number: // <-- Заменяем старый case float64
			node.Kind = KindNumber
			valStr := v.String()

			// Эвристика: если в строке числа есть точка или экспонента 'e'/'E', это float
			if strings.ContainsAny(valStr, ".eE") {
				if f, err := v.Float64(); err == nil {
					node.Value = strconv.FormatFloat(f, 'g', -1, 64)
				} else {
					node.Value = valStr
				}
			} else {
				// В противном случае парсим как чистый int64/uint64
				if i, err := v.Int64(); err == nil {
					node.Value = strconv.FormatInt(i, 10)
				} else {
					node.Value = valStr // fallback на сырую строку, если число гигантское
				}
			}
		case string:
			// Пытаемся раскрыть Base64
			if decoded, ok := h.tryDecodeBase64(v); ok {
				switch res := decoded.(type) {
				case string:
					// Успешно декодировали в Unicode текст
					node.Kind = KindString
					node.Value = res // Заменяем Base64 на чистый текст

				case []byte:
					// Это бинарные данные. Кодируем в красивую шестнадцатеричную строку
					node.Kind = KindString
					node.RawDisplay = true // Выводим без кавычек
					node.IsHex = true

					maxLen := len(res)
					truncated := false
					if maxLen > h.hexLimit && h.hexLimit > 0 {
						maxLen = h.hexLimit
						truncated = true
					}

					// Превращаем байты в hex-строку вида 0x7f07cea5...
					hexStr := hex.EncodeToString(res[:maxLen])

					if truncated {
						node.Value = fmt.Sprintf("%s... (%d bytes)", hexStr, len(res))
					} else {
						node.Value = hexStr
					}
				}
			} else {
				// Обычная строка, оставляем как есть
				node.Kind = KindString
				node.Value = v
			}
		default:
			node.Kind = KindNull
			node.Value = "null"
		}
		return node
	}

	node := &TreeNode{Key: key}
	switch delim {
	case '{':
		node.Kind = KindGroup
		for dec.More() {
			// Читаем ключ объекта — json.Decoder гарантирует исходный порядок!
			kToken, _ := dec.Token()
			k := kToken.(string)

			// Рекурсивно парсим значение для этого ключа
			node.Children = append(node.Children, h.parseJSONToken(k, dec))
		}
		dec.Token() // Закрываем '}'
	case '[':
		node.Kind = KindArray
		i := 0
		for dec.More() {
			node.Children = append(node.Children, h.parseJSONToken("["+strconv.Itoa(i)+"]", dec))
			i++
		}
		dec.Token() // Закрываем ']'
	}
	return node
}

func (h *SlogPrettyRenderer) renderIRTree(buf []byte, nodes []*TreeNode, states []bool, inErrorZone bool) []byte {
	count := len(nodes)
	for i, node := range nodes {
		isLast := i == count-1

		if node.Kind == KindStackTrace {
			buf = h.appendFormattedStackTrace(buf, node.Key, node.Value, states, isLast)
			continue
		}

		buf = append(buf, '\n')
		buf = append(buf, h.color.link...)
		for _, isParentLast := range states {
			if isParentLast {
				buf = append(buf, "   "...)
			} else {
				buf = append(buf, "│  "...)
			}
		}
		if isLast {
			buf = append(buf, "└── "...)
		} else {
			buf = append(buf, "├── "...)
		}
		buf = append(buf, h.color.reset...)

		// 1. Покраска ключа на основе точного семантического типа
		switch node.Kind {
		case KindErrorNode, KindErrorText:
			buf = append(buf, h.color.errkey...) // Инфраструктура ошибок (err, @text, @context)
		case KindLocation:
			buf = append(buf, h.color.loc...)
		case KindInlineArray:
			// Печатаем как примитив, но строго без кавычек
			buf = append(buf, h.color.ctx...) // используем цвет текста
			buf = append(buf, node.Value...)
		default:
			buf = append(buf, h.color.key...) // Обычные пользовательские ключи
		}
		buf = append(buf, node.Key...)
		buf = append(buf, h.color.reset...)

		isGroupType := node.Kind == KindGroup || node.Kind == KindArray || node.Kind == KindErrorNode
		if isGroupType {
			buf = h.renderIRTree(buf, node.Children, append(states, isLast), inErrorZone)
		} else {
			buf = append(buf, h.color.stdots...)
			buf = append(buf, ": "...)
			buf = append(buf, h.color.reset...)

			// 2. Покраска значения
			switch node.Kind {
			case KindLocation:
				buf = append(buf, h.color.loc...)
				buf = append(buf, node.Value...)
			case KindErrorText:
				buf = append(buf, h.color.error...) // Само тело ошибки горит красным
				buf = append(buf, node.Value...)
			case KindString:
				if node.IsHex {
					// 1. Печатаем приглушенный префикс "hex("
					buf = append(buf, h.color.stdots...)
					buf = append(buf, "hex("...)
					buf = append(buf, h.color.reset...)

					// 2. Печатаем само значение (его по-прежнему будет удобно выделять даблкликом!)
					buf = append(buf, h.color.ctx...)
					buf = append(buf, node.Value...)
					buf = append(buf, h.color.reset...)

					// 3. Печатаем приглушенную закрывающую скобку ")"
					buf = append(buf, h.color.stdots...)
					buf = append(buf, ")"...)
					buf = append(buf, h.color.reset...)
				} else {
					// Старая логика для обычных строк
					buf = append(buf, h.color.ctx...)
					if node.RawDisplay || strings.HasPrefix(node.Key, "[") {
						buf = append(buf, node.Value...)
					} else {
						buf = append(buf, '"')
						buf = append(buf, node.Value...)
						buf = append(buf, '"')
					}
					buf = append(buf, h.color.reset...)
				}
			case KindNull:
				buf = append(buf, h.color.trace...)
				buf = append(buf, node.Value...)
			case KindBool, KindNumber:
				buf = append(buf, h.color.debug...)
				buf = append(buf, node.Value...)
			default:
				buf = append(buf, h.color.ctx...)
				buf = append(buf, node.Value...)
			}
			buf = append(buf, h.color.reset...)
		}
	}
	return buf
}

func (h *SlogPrettyRenderer) appendFormattedStackTrace(buf []byte, key, stackStr string, states []bool, isCurrentLast bool) []byte {
	buf = append(buf, '\n')
	buf = append(buf, h.color.link...)
	for _, isParentLast := range states {
		if isParentLast {
			buf = append(buf, "   "...)
		} else {
			buf = append(buf, "│  "...)
		}
	}
	if isCurrentLast {
		buf = append(buf, "└── "...)
	} else {
		buf = append(buf, "├── "...)
	}
	buf = append(buf, h.color.reset...)

	buf = append(buf, h.color.errkey...)
	buf = append(buf, key...)
	buf = append(buf, h.color.reset...)
	buf = append(buf, h.color.stdots...)
	buf = append(buf, ": "...)
	buf = append(buf, h.color.reset...)

	fullStates := append(states, isCurrentLast)
	lines := strings.Split(stackStr, "\n")

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "goroutine ") {
			buf = h.appendStackLineIndent(buf, fullStates)
			buf = append(buf, h.color.panic...)
			buf = append(buf, ' ')
			buf = append(buf, line...)
			buf = append(buf, ' ')
			buf = append(buf, h.color.reset...)
			continue
		}
		if i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if strings.Contains(nextLine, ".go:") || strings.Contains(nextLine, "s:") {
				funcName := line
				locInfo := nextLine
				if idx := strings.LastIndex(locInfo, " "); idx != -1 {
					locInfo = locInfo[:idx]
				}
				buf = h.appendStackLineIndent(buf, fullStates)
				buf = append(buf, h.color.key...)
				buf = append(buf, funcName...)
				buf = append(buf, h.color.reset...)
				buf = append(buf, h.color.stdots...)
				buf = append(buf, " -> "...)
				buf = append(buf, h.color.reset...)
				buf = append(buf, h.color.loc...)
				buf = append(buf, locInfo...)
				buf = append(buf, h.color.reset...)
				i++
				continue
			}
		}
		buf = h.appendStackLineIndent(buf, fullStates)
		buf = append(buf, h.color.sttext...)
		buf = append(buf, line...)
		buf = append(buf, h.color.reset...)
	}
	return buf
}

func (h *SlogPrettyRenderer) appendStackLineIndent(buf []byte, fullStates []bool) []byte {
	buf = append(buf, '\n')
	buf = append(buf, h.color.link...)
	for _, isLast := range fullStates {
		if isLast {
			buf = append(buf, "   "...)
		} else {
			buf = append(buf, "│  "...)
		}
	}
	buf = append(buf, "   "...)
	buf = append(buf, h.color.reset...)
	return buf
}

// tryDecodeBase64 пытается разобрать строку.
// Если это валидный UTF-8 текст — возвращает его.
// Если бинарник — возвращает слайс байт.
// Если не Base64 — возвращает nil.
func (h *SlogPrettyRenderer) tryDecodeBase64(s string) (any, bool) {
	// Исключаем слишком короткие строки, чтобы избежать ложных срабатываний на обычных словах
	if len(s) < 4 {
		return nil, false
	}

	// Проверяем и декодируем Base64 (работаем со стандартным и URL-safe алфавитами)
	encoding := base64.StdEncoding
	if strings.ContainsAny(s, "-_") {
		encoding = base64.URLEncoding
	}

	decoded, err := encoding.DecodeString(s)
	if err != nil {
		return nil, false
	}

	// Проверяем, является ли результат валидной Unicode строкой
	if utf8.Valid(decoded) {
		// Дополнительный фильтр: проверяем, что это печатаемые символы, а не бинарный мусор, случайно совпавший с UTF-8
		isPrintable := true
		for _, r := range string(decoded) {
			if r < 32 && r != '\n' && r != '\r' && r != '\t' {
				isPrintable = false
				break
			}
		}
		if isPrintable {
			return string(decoded), true
		}
	}

	// Если не текст — возвращаем как бинарный слайс
	return decoded, true
}

func appendRawSlogValue(buf []byte, v slog.Value) []byte {
	switch v.Kind() {
	case slog.KindString:
		return append(buf, v.String()...)
	case slog.KindInt64:
		return strconv.AppendInt(buf, v.Int64(), 10)
	case slog.KindUint64:
		return strconv.AppendUint(buf, v.Uint64(), 10)
	case slog.KindFloat64:
		return strconv.AppendFloat(buf, v.Float64(), 'g', -1, 64)
	case slog.KindBool:
		return strconv.AppendBool(buf, v.Bool())
	case slog.KindTime:
		return v.Time().AppendFormat(buf, "2006-01-02 15:04:05.000")
	case slog.KindDuration:
		return append(buf, v.Duration().String()...)
	default:
		return append(buf, fmt.Sprint(v.Any())...)
	}
}

func (h *SlogPrettyRenderer) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogPrettyRenderer{
		opts:     h.opts,
		dst:      h.dst,
		color:    h.color,
		preAttrs: append(append([]slog.Attr{}, h.preAttrs...), attrs...),
	}
}

func (h *SlogPrettyRenderer) WithGroup(name string) slog.Handler { return h }

type forceTreeMarker struct{}

func (forceTreeMarker) LogValue() slog.Value {
	return slog.GroupValue(slog.String("__slog_force_tree__", ""))
}
func ForceTree() slog.Attr { return slog.Any("", forceTreeMarker{}) }
