package os

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	// "errors"

	"io/ioutil"
	"strings"

	"github.com/Felamande/god/lib/jsvm"
	"github.com/Felamande/otto"
)

func init() {
	p := jsvm.Module("os")
	p.Extend("readFile", readFile)
	p.Extend("readFileAsync", readFileAsync)
	p.Extend("system", system)
	p.Extend("getwd", getwd)
	p.Extend("wdName", wdName)
	p.Extend("writeFile", writeFile)
	p.Extend("rename", rename)
	p.Extend("exec", execfn)
	p.Extend("setenv", setenv)
}

func readFile(call otto.FunctionCall) otto.Value {

	file, err := call.Argument(0).ToString()
	errCb := call.Argument(1)
	if err != nil {
		return jsvm.Callback(errCb, "invalid file name "+file)
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}
	return jsvm.StringValue(string(b))

}

func readFileAsync(call otto.FunctionCall) otto.Value {
	fileArg := call.Argument(0)
	contentCb := call.Argument(1)
	errCb := call.Argument(2)

	file := fileArg.String()
	if !fileArg.IsString() {
		return jsvm.Callback(errCb, "invalid file name "+file)
	}

	f, err := os.Open(file)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	go func() {
		defer f.Close()
		buf := bytes.NewBuffer([]byte(""))
		io.Copy(buf, f)
		jsvm.Callback(contentCb, buf.String())
		jsvm.Callback(errCb, nil)
	}()
	return otto.UndefinedValue()

}

func getwd(call otto.FunctionCall) otto.Value {
	dir, err := os.Getwd()
	if err != nil {
		return otto.UndefinedValue()
	}
	v, _ := otto.ToValue(dir)

	return v
}

func system(call otto.FunctionCall) otto.Value {
	// fmt.Println("here")
	var cmdList []string

	arg0 := call.Argument(0)
	errCb := call.Argument(1)
	outPutCb := call.Argument(2)
	iarg0, err := arg0.Export()

	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	switch arg0t := iarg0.(type) {
	case []string:
		cmdList = arg0t
	case string:
		cmdList = strings.Split(arg0t, " ")
	default:
		return jsvm.Callback(errCb, "invalid commandline")
	}

	cmdLen := len(cmdList)
	if cmdLen == 0 {
		return jsvm.Callback(errCb, "no cmd")
	}
	c := exec.Command(cmdList[0], cmdList[1:]...)

	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	c.Stdout = stdout
	c.Stderr = stderr
	err = c.Run()
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}
	jsvm.Callback(errCb, stderr.String())
	jsvm.Callback(outPutCb, stdout.String())

	return otto.UndefinedValue()
}
func execfn(call otto.FunctionCall) otto.Value {
	var cmdList []string

	arg0 := call.Argument(0)
	errCb := call.Argument(1)
	iarg0, err := arg0.Export()

	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	switch arg0t := iarg0.(type) {
	case []string:
		cmdList = arg0t
	case string:
		cmdList = strings.Split(arg0t, " ")
	default:
		return jsvm.Callback(errCb, "invalid commandline")
	}

	cmdLen := len(cmdList)
	if cmdLen == 0 {
		return jsvm.Callback(errCb, "no cmd")
	}
	c := exec.Command(cmdList[0], cmdList[1:]...)
	c.Stdout = os.Stdout
	go func() {
		err = c.Run()
		if err != nil {
			jsvm.Callback(errCb, err.Error())
		}
	}()

	// os.Stdout.Write(b)
	return otto.UndefinedValue()
}

//writeFile function writeFile(file,content,flag,errCb,formatter)
//flag a:append, c:create, t:truncate, l:newline,s:sync
//formatter function, format content to custom string.
//errCb error callback
func writeFile(call otto.FunctionCall) otto.Value {
	fileArg := call.Argument(0)
	content := call.Argument(1)
	flagArg := call.Argument(2)
	errCb := call.Argument(3) //error callback
	formatterCb := call.Argument(4)
	// callbackArg := call.Argument(2)
	if !fileArg.IsString() {
		return jsvm.Callback(errCb, "invalid file name "+fileArg.String())
	}
	if content.IsUndefined() {
		return jsvm.Callback(errCb, "content is not provided.")
	}

	var flagStr string
	if flagArg.IsString() {
		flagStr = flagArg.String()
	}
	var flag = os.O_WRONLY
	var newline string
	for _, m := range flagStr {
		switch m {
		case 'a':
			flag |= os.O_APPEND
		case 'c':
			flag |= os.O_CREATE
		case 't':
			flag |= os.O_TRUNC
		case 's':
			flag |= os.O_SYNC
		case 'l':
			newline = "\n"
		}
	}
	go func() {
		f, err := os.OpenFile(fileArg.String(), flag, 0777)
		if err != nil {
			jsvm.Callback(errCb, err.Error())
			return
		}
		defer f.Close()
		formatted, err := jsvm.CbGetValue(formatterCb, content)
		if err != nil {
			jsvm.Callback(errCb, err.Error())
			return
		}
		bytes := []byte(formatted + newline)

		_, err = f.Write(bytes)
		if err != nil {
			jsvm.Callback(errCb, err.Error())
			return
		}
		jsvm.Callback(errCb, nil)
	}()
	return otto.UndefinedValue()

}

func wdName(call otto.FunctionCall) otto.Value {
	d, _ := os.Getwd()
	return jsvm.StringValue(filepath.Base(d))
}

func rename(call otto.FunctionCall) otto.Value {
	errCb := call.Argument(3)
	err := os.Rename(call.Argument(0).String(), call.Argument(1).String())
	if err != nil {
		jsvm.Callback(errCb, err.Error())
	}
	return otto.UndefinedValue()

}

func setenv(call otto.FunctionCall) otto.Value {
	objv := call.Argument(0)
	errCb := call.Argument(1)
	if !objv.IsObject() {
		return jsvm.Callback(errCb, "invalid object.")
	}

	obj := objv.Object()
	obj.ForEach(func(key string) {
		v, _ := obj.Get(key)
		os.Setenv(key, v.String())
	})
	return otto.UndefinedValue()

}
