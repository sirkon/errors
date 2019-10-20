package errors

import (
	"github.com/rs/zerolog"
)

// ContextReportBuilder абстракция занесения сущностей сохранённых ContextBuilder в контекст, например, логгера
type ContextReportBuilder interface {
	Int(name string, value int)
	Int8(name string, value int8)
	Int16(name string, value int16)
	Int32(name string, value int32)
	Int64(name string, value int64)
	Uint(name string, value uint)
	Uint8(name string, value uint8)
	Uint16(name string, value uint16)
	Uint32(name string, value uint32)
	Uint64(name string, value uint64)
	Float32(name string, value float32)
	Float64(name string, value float64)
	String(name string, value string)
	Strings(name string, values []string)
	Any(name string, value interface{})
	Err(err error)
}

var _ ContextReportBuilder = &zerologReportBuilder{}

// реализация contextReportBuilder-а для zerolog-а
type zerologReportBuilder struct {
	ctx zerolog.Context
}

// Logger получение логгера с данным контекстом
func (g *zerologReportBuilder) Logger(logger zerolog.Logger) zerolog.Logger {
	return g.ctx.Logger()
}

func (g *zerologReportBuilder) Int(name string, value int) {
	g.ctx = g.ctx.Int(name, value)
}

func (g *zerologReportBuilder) Int8(name string, value int8) {
	g.ctx = g.ctx.Int8(name, value)
}

func (g *zerologReportBuilder) Int16(name string, value int16) {
	g.ctx = g.ctx.Int16(name, value)
}

func (g *zerologReportBuilder) Int32(name string, value int32) {
	g.ctx = g.ctx.Int32(name, value)
}

func (g *zerologReportBuilder) Int64(name string, value int64) {
	g.ctx = g.ctx.Int64(name, value)
}

func (g *zerologReportBuilder) Uint(name string, value uint) {
	g.ctx = g.ctx.Uint(name, value)
}

func (g *zerologReportBuilder) Uint8(name string, value uint8) {
	g.ctx = g.ctx.Uint8(name, value)
}

func (g *zerologReportBuilder) Uint16(name string, value uint16) {
	g.ctx = g.ctx.Uint16(name, value)
}

func (g *zerologReportBuilder) Uint32(name string, value uint32) {
	g.ctx = g.ctx.Uint32(name, value)
}

func (g *zerologReportBuilder) Uint64(name string, value uint64) {
	g.ctx = g.ctx.Uint64(name, value)
}

func (g *zerologReportBuilder) Float32(name string, value float32) {
	g.ctx = g.ctx.Float32(name, value)
}

func (g *zerologReportBuilder) Float64(name string, value float64) {
	g.ctx = g.ctx.Float64(name, value)
}

func (g *zerologReportBuilder) String(name string, value string) {
	g.ctx = g.ctx.Str(name, value)
}

func (g *zerologReportBuilder) Strings(name string, values []string) {
	g.ctx = g.ctx.Strs(name, values)
}

func (g *zerologReportBuilder) Any(name string, value interface{}) {
	g.ctx = g.ctx.Interface(name, value)
}

func (g *zerologReportBuilder) Err(err error) {
	g.ctx = g.ctx.Err(err)
}
