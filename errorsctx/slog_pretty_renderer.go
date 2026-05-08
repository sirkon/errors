package errorsctx

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/sirkon/errors"
)

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

	// Собираем и трансформируем атрибуты на лету
	allAttrs := make([]slog.Attr, 0, len(h.preAttrs)+r.NumAttrs())
	allAttrs = append(allAttrs, h.preAttrs...)

	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, h.transformAttr(a))
		return true
	})

	hasGroups := false
	hasMultilineString := false

	for _, a := range allAttrs {
		if a.Value.Kind() == slog.KindGroup {
			hasGroups = true
			break
		}
		// Если это строка и в ней есть перенос, заставляем хендлер рисовать дерево
		if a.Value.Kind() == slog.KindString {
			if strings.Contains(a.Value.String(), "\n") {
				hasMultilineString = true
			}
		}
	}

	// Сценарий 1: Мало контекста, нет групп и многострочных строк -> Красивый раскрашенный JSON
	if !hasGroups && !hasMultilineString && len(allAttrs) <= 3 {
		if len(allAttrs) == 0 {
			return nil
		}

		// Добавляем разделительный пробел после сообщения лога
		buf = append(buf, ' ')

		// Открывающая фигурная скобка (красится цветом stdots или link)
		buf = append(buf, h.color.stdots...)
		buf = append(buf, '{')
		buf = append(buf, h.color.reset...)

		for i, a := range allAttrs {
			if i > 0 {
				// Запятая и пробел между парами ключ-значение
				buf = append(buf, h.color.stdots...)
				buf = append(buf, ", "...)
				buf = append(buf, h.color.reset...)
			}

			// 1. Подсветка ключа в JSON
			buf = append(buf, h.color.key...)
			buf = append(buf, '"')
			buf = append(buf, a.Key...)
			buf = append(buf, '"')
			buf = append(buf, h.color.reset...)

			// Двоеточие с пробелом
			buf = append(buf, h.color.stdots...)
			buf = append(buf, ": "...)
			buf = append(buf, h.color.reset...)

			// 2. Подсветка значения в JSON (используем профильный цвет ctx)
			buf = append(buf, h.color.ctx...)
			if a.Value.Kind() == slog.KindString {
				// Строки в JSON по стандарту оборачиваем в кавычки
				buf = append(buf, '"')
				buf = appendSlogValue(buf, a.Value)
				buf = append(buf, '"')
			} else {
				buf = appendSlogValue(buf, a.Value)
			}
			buf = append(buf, h.color.reset...)
		}

		// Закрывающая фигурная скобка
		buf = append(buf, h.color.stdots...)
		buf = append(buf, '}')
		buf = append(buf, h.color.reset...)

		return nil
	}

	// Сценарий 2: Вывод дерева (если есть группы, много полей или многострочный стек)
	buf = h.printTreeList(buf, allAttrs, []bool{}, false)
	return nil
}

// transformAttr перехватывает ошибки sirkon/errors и раскладывает их в типизированные slog-группы
func (h *SlogPrettyRenderer) transformAttr(a slog.Attr) slog.Attr {
	e, ok := a.Value.Any().(error)
	if !ok {
		return a
	}
	err, ok := e.(*errors.Error)
	if !ok {
		err, ok = errors.AsType[*errors.Error](e)
		if !ok {
			return a
		}
	}

	if a.Key == "" || a.Key == "!BADKEY" {
		a.Key = "err"
	}

	return slog.GroupAttrs(
		a.Key,
		slog.String("@text", e.Error()),
		slog.GroupAttrs("@context", errors.SLogTreeContext(err)...),
	)
}

func (h *SlogPrettyRenderer) printTreeList(buf []byte, attrs []slog.Attr, parentStates []bool, inErrorZone bool) []byte {
	count := len(attrs)
	for i, a := range attrs {
		isLast := i == count-1
		buf = h.printTreeAttr(buf, a, isLast, parentStates, inErrorZone)
	}
	return buf
}

func (h *SlogPrettyRenderer) printTreeAttr(
	buf []byte,
	a slog.Attr,
	isLast bool,
	parentStates []bool,
	inErrorZone bool,
) []byte {
	buf = append(buf, '\n')

	// Отрисовка направляющих линий дерева (палочек)
	buf = append(buf, h.color.link...)
	for _, isParentLast := range parentStates {
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

	isGroup := a.Value.Kind() == slog.KindGroup
	currentElementIsErrorNode := keyIsErrorContext(a.Key, isGroup)
	isCurrentErrorZone := inErrorZone || currentElementIsErrorNode

	// Перехват стектрейса (по ключу stacktrace или stack)
	if (a.Key == "stacktrace" || a.Key == "stack") && a.Value.Kind() == slog.KindString {
		buf = h.appendFormattedStackTrace(buf, a.Key, a.Value.String(), parentStates, isLast)
		return buf
	}

	// 1. Подсветка ключа
	buf = append(buf, h.selectKeyColor(a.Key, isGroup)...)
	buf = append(buf, a.Key...)
	buf = append(buf, h.color.reset...)

	if isGroup {
		nextStates := append(parentStates, isLast)
		buf = h.printTreeList(buf, a.Value.Group(), nextStates, isCurrentErrorZone)
	} else {
		buf = append(buf, h.color.stdots...)
		buf = append(buf, ": "...)
		buf = append(buf, h.color.reset...)

		// 2. Подсветка значений
		switch a.Key {
		case "@location":
			buf = append(buf, h.color.loc...)
			buf = appendSlogValue(buf, a.Value)
		case "@text":
			buf = append(buf, h.color.error...)
			buf = appendSlogValue(buf, a.Value)
		default:
			buf = append(buf, h.color.ctx...)
			buf = appendSlogValue(buf, a.Value)
		}
		buf = append(buf, h.color.reset...)
	}
	return buf
}

func (h *SlogPrettyRenderer) selectKeyColor(key string, isGroup bool) string {
	if key == "@location" {
		return h.color.loc
	}
	if key == "@text" {
		return h.color.errkey // Ключ для текста ошибки выделен
	}

	// Для @context, NEW:, WRAP:, CTX и err используем стандартный key,
	// чтобы они не перегружали терминал красным цветом
	return h.color.key
}

// Точное определение инфраструктурных ключей пакета ошибок
func keyIsErrorContext(key string, isGroup bool) bool {
	if key == "err" || key == "@text" || key == "@context" || strings.HasPrefix(key, "err-") {
		return true
	}
	// Если это группа и она завершается двоеточием (NEW:, WRAP:), это слой вашего sirkon/errors
	if isGroup && strings.HasSuffix(key, ":") {
		return true
	}
	if key == "CTX" {
		return true
	}
	return false
}

func appendSlogValue(buf []byte, v slog.Value) []byte {
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
		b, err := json.Marshal(v.Any())
		if err != nil {
			return append(buf, err.Error()...)
		}
		return append(buf, b...)
	}
}

func (h *SlogPrettyRenderer) appendFormattedStackTrace(buf []byte, key, stackStr string, parentStates []bool, isCurrentLast bool) []byte {
	// Подсвечиваем заголовок stacktrace
	buf = append(buf, h.color.errkey...)
	buf = append(buf, key...)
	buf = append(buf, h.color.reset...)
	buf = append(buf, h.color.stdots...)
	buf = append(buf, ": "...)
	buf = append(buf, h.color.reset...)

	// Создаем полный стек состояний, включая состояние текущего узла stack
	fullStates := append(parentStates, isCurrentLast)

	lines := strings.Split(stackStr, "\n")

	// Перебираем строки стека парами: Функция -> Локация
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Вывод информации о горутине
		if strings.HasPrefix(line, "goroutine ") {
			buf = h.appendStackLineIndent(buf, fullStates)
			buf = append(buf, h.color.panic...)
			buf = append(buf, ' ')
			buf = append(buf, line...)
			buf = append(buf, ' ')
			buf = append(buf, h.color.reset...)
			continue
		}

		// Обработка пары Функция + Файл
		if i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if strings.Contains(nextLine, ".go:") || strings.Contains(nextLine, "s:") {
				funcName := line
				locInfo := nextLine

				if idx := strings.LastIndex(locInfo, " "); idx != -1 {
					locInfo = locInfo[:idx]
				}

				// Печать кадра стека
				buf = h.appendStackLineIndent(buf, fullStates)

				// Функция
				buf = append(buf, h.color.key...)
				buf = append(buf, funcName...)
				buf = append(buf, h.color.reset...)

				buf = append(buf, h.color.stdots...)
				buf = append(buf, " -> "...)
				buf = append(buf, h.color.reset...)

				// Путь
				buf = append(buf, h.color.loc...)
				buf = append(buf, locInfo...)
				buf = append(buf, h.color.reset...)

				i++
				continue
			}
		}

		// Фолбек
		buf = h.appendStackLineIndent(buf, fullStates)
		buf = append(buf, h.color.sttext...)
		buf = append(buf, line...)
		buf = append(buf, h.color.reset...)
	}

	return buf
}

// Теперь этот метод опирается на fullStates, где последний элемент — состояние самого узла stack
func (h *SlogPrettyRenderer) appendStackLineIndent(buf []byte, fullStates []bool) []byte {
	buf = append(buf, '\n')
	buf = append(buf, h.color.link...)

	// Отрисовываем все уровни отступов, включая уровень текущего стектрейса
	for _, isLast := range fullStates {
		if isLast {
			buf = append(buf, "   "...) // Если уровень закрылся, линии нет
		} else {
			buf = append(buf, "│  "...) // Если уровень продолжается, рисуем линию
		}
	}

	// Добавляем небольшой фиксированный сдвиг вправо для красоты текста
	buf = append(buf, "   "...)
	buf = append(buf, h.color.reset...)
	return buf
}

// Вспомогательный хелпер для сохранения вертикальных линий дерева внутри стектрейса
func (h *SlogPrettyRenderer) appendTreeIndent(buf []byte, parentStates []bool, isLast bool) []byte {
	buf = append(buf, '\n')
	buf = append(buf, h.color.link...)
	for _, isParentLast := range parentStates {
		if isParentLast {
			buf = append(buf, "   "...)
		} else {
			buf = append(buf, "│  "...)
		}
	}
	// Внутри стектрейса все дочерние элементы идут как продолжающиеся ветки родительского узла ошибки
	buf = append(buf, "│  ├── "...)
	buf = append(buf, h.color.reset...)
	return buf
}

func (h *SlogPrettyRenderer) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Предобработка атрибутов при ветвлении логгера через .With()
	transformed := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		transformed[i] = h.transformAttr(a)
	}
	return &SlogPrettyRenderer{
		opts:     h.opts,
		dst:      h.dst,
		color:    h.color,
		preAttrs: append(append([]slog.Attr{}, h.preAttrs...), transformed...),
	}
}

func (h *SlogPrettyRenderer) WithGroup(name string) slog.Handler {
	// Метод WithGroup можно оставить пустым, если структура плоская на верхнем уровне
	return h
}
