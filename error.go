package errors

import (
	"errors"
	"log/slog"
	"strings"
)

type Error struct {
	attrs []errorAttr
}

type errorAttr struct {
	key   string
	value slog.Value
	kind  errorAttrKind
}

func (e *Error) Error() string {
	links := make([]string, 0, 2*errorContextLengthPrediction)

	for _, attr := range e.attrs {
		switch attr.kind {
		case errorAttrKindNew:
			links = append(links, attr.key)
		case errorAttrKindWrap:
			links = append(links, attr.key)
		case errorAttrKindOutterWrap:
			links = append(links, attr.value.Any().(error).Error())
			links = append(links, attr.key)
		case errorAttrKindOutterJust:
			links = append(links, attr.value.Any().(error).Error())
		default:
			continue
		}
	}

	var msg strings.Builder
	for i := len(links) - 1; i >= 0; i-- {
		if i != len(links)-1 {
			msg.WriteString(": ")
		}
		msg.WriteString(links[i])
	}

	return msg.String()
}

// Is implements support for [Is]. Remember, *Error instances
// are ephemeral and cannot be targeted.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	if _, ok := target.(*Error); ok {
		return false
	}

	for _, attr := range e.attrs {
		switch attr.kind {
		case errorAttrKindOutterWrap, errorAttrKindOutterJust:
			if errors.Is(attr.value.Any().(error), target) {
				return true
			}
		default:
			continue
		}
	}

	return false
}

// As implements support for [Is]. Remember, *Error instances
// are ephemeral and cannot be targeted.
func (e *Error) As(target any) bool {
	if target == nil {
		return false
	}

	switch v := target.(type) {
	case **Error:
		return false
	case **errorContextDeliverer:
		(*v) = &errorContextDeliverer{
			tgt: e,
		}
		return true
	}

	for _, attr := range e.attrs {
		switch attr.kind {
		case errorAttrKindOutterWrap, errorAttrKindOutterJust:
			if errors.As(attr.value.Any().(error), target) {
				return true
			}
		default:
			continue
		}
	}

	return false
}

type errorAttrKind int8

const (
	errorAttrKindInvalid errorAttrKind = iota
	errorAttrKindNew
	errorAttrKindWrap
	errorAttrKindOutterWrap
	errorAttrKindJust
	errorAttrKindOutterJust
	errorAttrKindLoc
	errorAttrKindBool
	errorAttrKindI64
	errorAttrKindU64
	errorAttrKindF64
	errorAttrKindStr
	errorAttrKindAny
)
