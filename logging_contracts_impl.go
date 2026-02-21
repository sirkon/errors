package errors

type errorContextDeliverer struct {
	tgt *Error
}

func (e *errorContextDeliverer) Deliver(cons ErrorContextConsumer) {
	var layer ErrorContextBuilder

	for _, attr := range e.tgt.attrs {
		switch attr.kind {
		case errorAttrKindNew:
			finalizeErrContextBuilder(layer)
			layer = cons.New(attr.key)
		case errorAttrKindWrap:
			finalizeErrContextBuilder(layer)
			layer = cons.Wrap(attr.key)
		case errorAttrKindOutterWrap:
			finalizeErrContextBuilder(layer)
			dlv := GetContextDeliverer(attr.value.Any().(error))
			if dlv != nil {
				dlv.Deliver(cons)
			}
			layer = cons.Wrap(attr.key)
		case errorAttrKindJust:
			finalizeErrContextBuilder(layer)
			layer = cons.Just()
		case errorAttrKindOutterJust:
			finalizeErrContextBuilder(layer)
			layer = cons.Just()
			dlv := GetContextDeliverer(attr.value.Any().(error))
			if dlv != nil {
				dlv.Deliver(cons)
			}
			layer = cons.Just()
		case errorAttrKindLoc:
			layer.Loc(attr.value.String())
		case errorAttrKindBool:
			layer.Bool(attr.key, attr.value.Bool())
		case errorAttrKindI64:
			layer.Int64(attr.key, attr.value.Int64())
		case errorAttrKindU64:
			layer.Uint64(attr.key, attr.value.Uint64())
		case errorAttrKindF64:
			layer.Flt64(attr.key, attr.value.Float64())
		case errorAttrKindStr:
			layer.Str(attr.key, attr.value.String())
		case errorAttrKindAny:
			layer.Any(attr.key, attr.value.Any())
		default:
			panic(Newf("invalid errorAttrKind value %d for key %q", attr.kind, attr.key))
		}
	}

	if layer != nil {
		layer.Finalize()
	}
}

func (e *errorContextDeliverer) Error() string { return "" }

func finalizeErrContextBuilder(layer ErrorContextBuilder) {
	if layer == nil {
		return
	}
	layer.Finalize()
}
