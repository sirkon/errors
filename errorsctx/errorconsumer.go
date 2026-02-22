package errorsctx

import (
	"fmt"
	"log/slog"

	"github.com/sirkon/errors"
)

// Consumer an implementation of [errors.ErrorContextConsumer] providing
// a tree view of the collected context grouped by layers of processing nodes.
type Consumer struct {
	Layers []Layer
}

func (c *Consumer) New(msg string) errors.ErrorContextBuilder {
	return &Layer{
		parent: c,
		Kind:   LayerKindNew,
		What:   msg,
		Pairs:  make([]slog.Attr, 0, errorContextLengthPrediction),
	}
}

func (c *Consumer) Wrap(msg string) errors.ErrorContextBuilder {
	return &Layer{
		parent: c,
		Kind:   LayerKindWrap,
		What:   msg,
		Pairs:  make([]slog.Attr, 0, errorContextLengthPrediction),
	}
}

func (c *Consumer) Just() errors.ErrorContextBuilder {
	return &Layer{
		parent: c,
		Kind:   LayerKindJust,
		Pairs:  make([]slog.Attr, 0, errorContextLengthPrediction),
	}
}

type Layer struct {
	parent *Consumer
	Kind   LayerKind
	What   string
	Pos    string
	Pairs  []slog.Attr
}

func (c Layer) String() string {
	if c.What == "" {
		return c.Kind.String()
	}

	return c.Kind.String() + ": " + c.What
}

func (c *Layer) Bool(name string, value bool) {
	c.Pairs = append(c.Pairs, slog.Bool(name, value))
}

func (c *Layer) Int64(name string, value int64) {
	c.Pairs = append(c.Pairs, slog.Int64(name, value))
}

func (c *Layer) Uint64(name string, value uint64) {
	c.Pairs = append(c.Pairs, slog.Uint64(name, value))
}

func (c *Layer) Flt64(name string, value float64) {
	c.Pairs = append(c.Pairs, slog.Float64(name, value))
}

func (c *Layer) Str(name string, value string) {
	c.Pairs = append(c.Pairs, slog.String(name, value))
}

func (c *Layer) Any(name string, value any) {
	c.Pairs = append(c.Pairs, slog.Any(name, value))
}

func (c *Layer) Loc(position string) {
	c.Pos = position
}

func (c *Layer) Finalize() {
	if c.parent == nil {
		return
	}
	c.parent.Layers = append(c.parent.Layers, *c)
	c.parent = nil
}

type LayerKind int

func (l LayerKind) String() string {
	switch l {
	case LayerKindNew:
		return "NEW"
	case LayerKindWrap:
		return "WRAP"
	case LayerKindJust:
		return "CTX"
	default:
		return fmt.Sprintf("InvalidLayerKind(%d)", l)
	}
}

const (
	layerKindInvalid LayerKind = iota
	LayerKindNew
	LayerKindWrap
	LayerKindJust
)
