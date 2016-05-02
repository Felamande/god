package localstorage

import (
	"github.com/Felamande/god/lib/jsvm"
	"github.com/robertkrimen/otto"
)

type Cloner interface {
	Clone() Storage
}

type Storage interface {
	Put(key []byte, value []byte) error
	Get(key []byte) (value []byte, err error)
}

func makePutFn(s Storage) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0).String()
		val := call.Argument(1).String()
		errCb := call.Argument(2)
		var err error
		if cloner, ok := s.(Cloner); ok {
			err = cloner.Clone().Put([]byte(key), []byte(val))
		} else {
			err = s.Put([]byte(key), []byte(val))
		}

		if err != nil {
			return jsvm.Callback(errCb, err.Error())
		}
		return otto.UndefinedValue()
	}

}
func makeGetFn(s Storage) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		key := call.Argument(0).String()
		errCb := call.Argument(1)
		var val []byte
		var err error
		if cloner, ok := s.(Cloner); ok {
			val, err = cloner.Clone().Get([]byte(key))
		} else {
			val, err = s.Get([]byte(key))
		}
		if err != nil {
			return jsvm.Callback(errCb, err.Error())
		}
		return jsvm.StringValue(string(val))
	}
}
