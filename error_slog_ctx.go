package errors

import (
	"log/slog"
)

// SLogTreeContext builds slog tree for logging on the place without
// [ErrorContextConsumer] indirect calls overhead.
func SLogTreeContext(err *Error) []slog.Attr {
	s := newSlogTreeContextState()
	s.feed(err.attrs)
	return s.stages
}

func SLogFlatContext(err *Error) []slog.Attr {
	s := newSlogFlatContextState()
	s.feed(err.attrs)
	s.ctx = append(s.ctx, slog.GroupAttrs("@locations", s.pos...))
	return s.ctx
}

type slogTreeContextState struct {
	stages       []slog.Attr
	stage        []slog.Attr
	hasPos       bool
	hasSomething bool
	name         string
}

func newSlogTreeContextState() *slogTreeContextState {
	return &slogTreeContextState{
		stages: make([]slog.Attr, 0, errorContextNoOfStages),
		stage:  make([]slog.Attr, 1, errorContextLengthPrediction),
	}
}

func (s *slogTreeContextState) feed(attrs []errorAttr) {
	for _, attr := range attrs {
		switch attr.kind {
		case errorAttrKindNew:
			s.closeStage()
			s.name = "NEW: " + attr.key

		case errorAttrKindWrap:
			s.closeStage()
			s.name = "WRAP: " + attr.key

		case errorAttrKindOutterWrap:
			// Рекурсивный спуск.
			if e, ok := attr.value.Any().(error); ok {
				nerr, ok := AsType[*Error](e)
				if ok {
					s.feed(nerr.attrs)
				}
			}

			s.closeStage()
			s.name = "WRAP: " + attr.key

		case errorAttrKindJust:
			s.closeStage()
			s.name = "CTX"

		case errorAttrKindOutterJust:
			// Рекурсивный спуск.
			if e, ok := attr.value.Any().(error); ok {
				nerr, ok := AsType[*Error](e)
				if ok {
					s.feed(nerr.attrs)
				}
			}

			s.closeStage()
			s.name = "CTX"

		case errorAttrKindLoc:
			s.stage[0] = slog.Attr{
				Key:   "@location",
				Value: attr.value,
			}
		default:
			s.stage = append(s.stage, slog.Attr{
				Key:   attr.key,
				Value: attr.value,
			})
		}
	}

	s.closeStage()
}

func (s *slogTreeContextState) closeStage() {
	if s.name == "" {
		return
	}

	s.stages = append(s.stages, slog.GroupAttrs(s.name, s.stage...))
	s.stage = s.stage[len(s.stage):]
	if cap(s.stage) > 0 {
		s.stage = s.stage[:1]
	} else {
		s.stage = make([]slog.Attr, 1, errorContextLengthPrediction)
	}
}

type slogFlatContextState struct {
	ctx  []slog.Attr
	pos  []slog.Attr
	name string
}

func newSlogFlatContextState() *slogFlatContextState {
	return &slogFlatContextState{
		ctx: make([]slog.Attr, 0, errorContextLengthPrediction),
		pos: make([]slog.Attr, 0, errorContextNoOfStages),
	}
}

func (s *slogFlatContextState) feed(attrs []errorAttr) {
	for _, attr := range attrs {
		switch attr.kind {
		case errorAttrKindNew:
			s.name = "NEW: " + attr.key
		case errorAttrKindWrap:
			s.name = "WRAP: " + attr.key
		case errorAttrKindOutterWrap:
			// Рекурсивный спуск.
			if e, ok := attr.value.Any().(error); ok {
				nerr, ok := AsType[*Error](e)
				if ok {
					s.feed(nerr.attrs)
				}
			}
			s.name = "WRAP: " + attr.key
		case errorAttrKindJust:
			s.name = "CTX"
		case errorAttrKindOutterJust:
			// Рекурсивный спуск.
			if e, ok := attr.value.Any().(error); ok {
				nerr, ok := AsType[*Error](e)
				if ok {
					s.feed(nerr.attrs)
				}
			}
			s.name = "CTX"
		case errorAttrKindLoc:
			s.pos = append(s.pos, slog.String(s.name, attr.value.String()))
		default:
			s.ctx = append(s.ctx, slog.Attr{
				Key:   attr.key,
				Value: attr.value,
			})
		}
	}
}
