package jsvm

import (
	"errors"
	"os"

	"github.com/robertkrimen/otto"
)

type Func otto.Value

func (f Func) Call(args ...interface{}) {
	Callback(otto.Value(f), args...)
}

type module struct {
	obj      *otto.Object
	methods  map[string]interface{}
	required bool
	initfn   func()
	deferfn  func()

	// register string
}

type builtin struct {
	obj *otto.Object
}

var modules map[string]*module
var builtins map[string]*builtin

var vm *otto.Otto

func init() {
	modules = make(map[string]*module)
	builtins = make(map[string]*builtin)
	vm = otto.New()
	vm.Set("require", require)

}

func Vm() *otto.Otto {
	return vm
}

func Builtin(obj string) *builtin {
	if b, exist := builtins[obj]; exist {
		return b
	}
	if _, exist := modules[obj]; exist {
		return nil
	}
	v, err := vm.Get(obj)
	if err != nil {
		return nil
	}
	return &builtin{v.Object()}

}

func (b *builtin) Extend(name string, fn interface{}) error {
	return b.obj.Set(name, fn)
}

func require(call otto.FunctionCall) otto.Value {
	mod := call.Argument(0).String()
	p, exist := modules[mod]
	if !exist {
		return otto.UndefinedValue()
	}

	if p.required {
		return p.obj.Value()
	}
	for name, method := range p.methods {
		p.obj.Set(name, method)
	}
	p.required = true
	if p.initfn != nil {
		p.initfn()
	}
	return p.obj.Value()

}

func Module(name string) *module {
	if m, exist := modules[name]; exist {
		return m

	}
	o, _ := vm.Object(`({})`)

	vm.Set(name, o)

	m := &module{obj: o, methods: make(map[string]interface{}), required: false}
	modules[name] = m
	return m
}

func (m *module) Init(f func()) {
	m.initfn = f
}

func (m *module) Defer(f func()) {
	m.deferfn = f
}

func (m *module) Extend(obj string, Func func(call otto.FunctionCall) otto.Value) {
	m.methods[obj] = Func
}
func (m *module) Obj() *otto.Object {
	return m.obj
}

func Run(src string) error {
	var script *otto.Script
	var err error
	if fi, _ := os.Stat(src); fi == nil {
		script, err = vm.Compile("", src)

	} else {

		script, err = vm.Compile(src, nil)
	}

	if err != nil {
		return err
	}

	_, err = vm.Run(script)
	return err
}

func StringValue(s string) otto.Value {
	v, _ := otto.ToValue(s)
	return v
}

func ErrorValue(err error) otto.Value {
	value, _ := otto.ToValue(err.Error())
	return value
}

func Callback(cb otto.Value, arg ...interface{}) otto.Value {
	if cb.Class() != "Function" {
		return otto.UndefinedValue()
	}
	cb.Call(cb, arg...)
	return otto.UndefinedValue()

}
func CbGetValue(cb otto.Value, arg otto.Value) (string, error) {
	if cb.IsUndefined() {
		return arg.String(), nil
	}
	if !cb.IsFunction() {
		return "", errors.New("invalid formatter")
	}
	v, err := cb.Call(cb, arg)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func ToObject(o O) otto.Value {
	oo, _ := vm.Object(`({})`)
	for name, value := range o {
		oo.Set(name, value)
	}
	return oo.Value()
}

type O map[string]interface{}

func CallDefer() {
	for _, m := range modules {
		m.deferfn()
	}
}
