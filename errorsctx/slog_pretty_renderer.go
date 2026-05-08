package errorsctx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/sirkon/errors"
)

// NodeKind определяет тип конечного значения для точечной раскраски
type NodeKind int

const (
	KindGroup NodeKind = iota
	KindArray
	KindString
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
	Key      string
	Value    string
	Kind     NodeKind
	Children []*TreeNode
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
}

func NewSlogPrettyRenderer(dst io.Writer, opts *slog.HandlerOptions, isDark bool) *SlogPrettyRenderer {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	var profile *prettyWriterColorProfile
	if isDark {
		profile = newPrettyWriterColorProfileDark()
	} else {
		profile = newPrettyWriterColorProfileLight()
	}
	return &SlogPrettyRenderer{
		opts:  *opts,
		dst:   dst,
		color: profile,
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

	if len(rawStr) > 1 && (rawStr[0] == '{' || rawStr[0] == '[') {
		if json.Unmarshal([]byte(rawStr), &anyObj) == nil {
			isJSON = true
		}
	} else if resolved.Kind() == slog.KindAny {
		if jsonBytes, err := json.Marshal(resolved.Any()); err == nil && len(jsonBytes) > 1 {
			if jsonBytes[0] == '{' || jsonBytes[0] == '[' {
				if json.Unmarshal(jsonBytes, &anyObj) == nil {
					isJSON = true
				}
			}
		}
	}

	if isJSON {
		return h.buildIRTreeFromAny(key, anyObj)
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

// buildIRTreeFromAny рекурсивно раскладывает any-объекты (мапы/слайсы из JSON) в IR-узлы
func (h *SlogPrettyRenderer) buildIRTreeFromAny(key string, obj any) *TreeNode {
	node := &TreeNode{Key: key}
	if obj == nil {
		node.Kind = KindNull
		node.Value = "null"
		return node
	}

	switch v := obj.(type) {
	case map[string]any:
		if len(v) == 0 {
			node.Kind = KindNull
			node.Value = "{}"
			return node
		}
		node.Kind = KindGroup
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			node.Children = append(node.Children, h.buildIRTreeFromAny(k, v[k]))
		}

	case []any:
		if len(v) == 0 {
			node.Kind = KindNull
			node.Value = "[]"
			return node
		}
		node.Kind = KindArray
		for i, val := range v {
			node.Children = append(node.Children, h.buildIRTreeFromAny("["+strconv.Itoa(i)+"]", val))
		}

	case string:
		if strings.Contains(v, ".go:") || key == "@location" {
			node.Kind = KindLocation
		} else {
			node.Kind = KindString
		}
		node.Value = v

	case bool:
		node.Kind = KindBool
		node.Value = strconv.FormatBool(v)

	case float64:
		node.Kind = KindNumber
		node.Value = strconv.FormatFloat(v, 'g', -1, 64)

	default:
		node.Kind = KindNumber
		node.Value = fmt.Sprint(v)
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
				buf = append(buf, h.color.ctx...)
				if !strings.HasPrefix(node.Key, "[") {
					buf = append(buf, '"')
					buf = append(buf, node.Value...)
					buf = append(buf, '"')
				} else {
					buf = append(buf, node.Value...)
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
