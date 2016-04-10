package hotkey

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Felamande/god/lib/kbevent"
	"github.com/Felamande/jsvm"
	"github.com/Felamande/otto"
)

var hkCb = make(map[string]jsvm.Func)

func init() {
	if m := jsvm.Module("hotkey"); m != nil {
		m.Extend("register", register)
	}
}

func register(call otto.FunctionCall) otto.Value {
	hk := call.Argument(0).String()
	cb := call.Argument(1)
	hkCb[hk] = jsvm.Func(cb)
	// fmt.Print(hk)
	// jsvm.Func(cb).Call()
	return otto.UndefinedValue()
}

func ApplyAll() error {
	runtime.LockOSThread()

	kb := kbevent.New()
	kb.Bind("ctrl+d", func() { os.Exit(0) })

	for hk, cb := range hkCb {
		err := kb.Bind(hk, makeFunc(cb))
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}

	return kb.Start()
}

func makeFunc(cb jsvm.Func) func() {
	return func() { cb.Call() }
}
