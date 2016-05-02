package localstorage

import (
	"fmt"

	"github.com/Felamande/god/lib/jsvm"
	"github.com/Felamande/otto"
)

var defaultType = "gkvlite"
var localStorageType map[string]Storage

func init() {
	if m := jsvm.Module("localstorage"); m != nil {
		localStorageType = make(map[string]Storage)
		localStorageType["bolt"] = &boltStorage{}
		localStorageType["gomap"] = &mapStorage{make(map[string]string)}
		localStorageType["leveldb"] = &lvdbStorage{}
		// localStorageType["qldb"] = &qlStorage{}
		localStorageType["gkvlite"] = &kvdbStorage{}
		m.Extend("put", put)
		m.Extend("get", get)
		m.Extend("use", use)
		m.Extend("setbackend", setbackend)
	}
}

func setbackend(call otto.FunctionCall) otto.Value {
	typ := call.Argument(0).String()
	defaultType = typ
	return otto.UndefinedValue()
}

func use(call otto.FunctionCall) otto.Value {
	typ := call.Argument(0).String()
	errCb := call.Argument(1)
	storage, exist := localStorageType[typ]
	if !exist {
		return jsvm.Callback(errCb, "unsupported backend "+typ)
	}
	return jsvm.ToObject(jsvm.O{"put": makePutFn(storage), "get": makeGetFn(storage)})
}

func put(call otto.FunctionCall) otto.Value {

	storage, ok := localStorageType[defaultType]
	if !ok {
		fmt.Println("unsupported backend", defaultType)
		return otto.UndefinedValue()
	}

	return makePutFn(storage)(call)
}

func get(call otto.FunctionCall) otto.Value {
	storage, ok := localStorageType[defaultType]
	if !ok {
		fmt.Println("unsupported backend", defaultType)
		return otto.UndefinedValue()
	}
	return makeGetFn(storage)(call)
}
