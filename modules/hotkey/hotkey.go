package hotkey

import (
	"github.com/Felamande/god/lib/jsvm"
	"github.com/Felamande/god/lib/kbevent"
	"github.com/Felamande/otto"
)

var kb = kbevent.New()

func init() {
	if m := jsvm.Module("hotkey"); m != nil {
		m.Extend("bind", bind)
	}
}

func bind(call otto.FunctionCall) otto.Value {
	hk := call.Argument(0).String()
	cb := call.Argument(1)
	errCb := call.Argument(2)
	err := kb.Bind(hk, func() { jsvm.Callback(cb) })
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	return otto.UndefinedValue()
}

func Bind(seq string, f func()) error {
	return kb.Bind(seq, f)
}

func ApplyAll() {
	kb.Start(nil)
}

func makeFunc(cb jsvm.Func) func() {
	return func() { cb.Call() }
}
