package path

import (

	// "errors"

	"path"
	"path/filepath"

	"github.com/Felamande/god/lib/jsvm"
	"github.com/robertkrimen/otto"
)

func init() {
	p := jsvm.Module("path")
	p.Extend("join", join)
	p.Extend("dir", dir)
}

func join(call otto.FunctionCall) otto.Value {
	var pathList []string

	for _, arg := range call.ArgumentList {
		v, _ := arg.Export()
		switch val := v.(type) {
		case []string:
			pathList = append(pathList, val...)
		case string:
			pathList = append(pathList, val)
		}
	}
	path := filepath.Join(pathList...)
	v, _ := otto.ToValue(filepath.ToSlash(path))
	return v
}

func dir(call otto.FunctionCall) otto.Value {
	p := call.Argument(0).String()
	p = filepath.ToSlash(p)
	if filepath.IsAbs(p) {
		return jsvm.StringValue(path.Dir(p))
	}
	var suffix string
	dir := path.Dir(p)
	if dir != "." && dir != ".." {
		suffix = "./"
	}
	return jsvm.StringValue(suffix + dir)
}
