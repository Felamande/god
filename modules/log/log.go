package log

import (
	"os"

	"github.com/Felamande/jsvm"
	"github.com/Felamande/otto"
	"github.com/qiniu/log"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[daemon]", log.LstdFlags|log.Llevel)
	if p := jsvm.Module("log"); p != nil {
		p.Extend("info", info)
		p.Extend("error", err)
		p.Extend("warn", warn)
	}
}

func info(call otto.FunctionCall) otto.Value {
	logger.Info(format(call)...)
	return otto.UndefinedValue()
}

func err(call otto.FunctionCall) otto.Value {
	logger.Error(format(call)...)
	return otto.UndefinedValue()
}
func warn(call otto.FunctionCall) otto.Value {
	logger.Warn(format(call)...)
	return otto.UndefinedValue()
}

func format(call otto.FunctionCall) (re []interface{}) {
	re = append(re, call.CallerLocation())
	for _, v := range call.ArgumentList {
		iarg, _ := v.Export()
		re = append(re, iarg)
	}
	return
}
