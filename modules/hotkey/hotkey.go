package hotkey

import (
	"os"

	"github.com/Felamande/god/lib/kbevent"
	"github.com/Felamande/jsvm"
	"github.com/Felamande/otto"
)

var kb = kbevent.New()

func init() {
	kb.Bind("ctrl+d", func() { os.Exit(0) })
	if m := jsvm.Module("hotkey"); m != nil {
		m.Extend("register", register)
	}
}

func register(call otto.FunctionCall) otto.Value {
	hk := call.Argument(0).String()
	cb := call.Argument(1)
	errCb := call.Argument(2)
	err := kb.Bind(hk, func() { jsvm.Callback(cb) })
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	return otto.UndefinedValue()
}

func ApplyAll() {
	kb.Start(nil)
}

func makeFunc(cb jsvm.Func) func() {
	return func() { cb.Call() }
}
