package gofunc

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Felamande/god/process"
	"github.com/Felamande/jsvm"
	"github.com/Felamande/otto"
)

var actionTime = make(map[string]time.Time)

func init() {
	// wd, _ = os.Getwd()
	if p := jsvm.Module("go"); p != nil {
		p.Extend("build", build)
		p.Extend("install", install)
		p.Extend("test", test)
		p.Extend("reload", reload)
	}
}

func build(call otto.FunctionCall) otto.Value {
	pkg := call.Argument(0).String()
	Args := call.Argument(1)
	errCb := call.Argument(2)

	abs, err := filepath.Abs(pkg)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	pkgName := filepath.Base(abs)
	err = exect(nil, "go", "build", "-o", getBinName(pkg, pkgName), Args, pkg)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}
	return otto.UndefinedValue()
}

func install(call otto.FunctionCall) otto.Value {
	if lastTime, ok := actionTime["install"]; ok {
		if time.Now().Sub(lastTime) < 2 {
			return otto.UndefinedValue()
		}
	}

	pkg := call.Argument(0).String()
	Args := call.Argument(1)
	errCb := call.Argument(2)

	err := exect(nil, "go", "install", Args, pkg)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}
	actionTime["install"] = time.Now()

	return otto.UndefinedValue()
}

func test(call otto.FunctionCall) otto.Value {
	if lastTime, ok := actionTime["test"]; ok {
		if time.Now().Sub(lastTime) < 2 {
			return otto.UndefinedValue()
		}
	}

	pkg := call.Argument(0).String()
	Args := call.Argument(1)
	errCb := call.Argument(2)

	err := exect(os.Stdout, "go", "test", Args, pkg)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	actionTime["test"] = time.Now()
	return otto.UndefinedValue()
}

func reload(call otto.FunctionCall) otto.Value {
	if lastTime, ok := actionTime["reload"]; ok {
		if time.Now().Sub(lastTime) < 2 {
			return otto.UndefinedValue()
		}
	}

	pkgPath := call.Argument(0).String()
	buildArgs := call.Argument(1)
	binArgs := call.Argument(2)
	errCb := call.Argument(3)

	abs, err := filepath.Abs(pkgPath)
	if err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	pkgName := filepath.Base(abs)

	tmpBin := getBinName(pkgPath, pkgName+"_tmp")

	if err = exect(nil, "go", "build", buildArgs, "-o", tmpBin, pkgPath); err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	if _, err = os.Stat(tmpBin); err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	if err = process.KillByName(pkgName); err != nil {
		jsvm.Callback(errCb, err.Error())
	}

	bin := getBinName(pkgPath, pkgName)
	if err = os.Remove(bin); err != nil {
		jsvm.Callback(errCb, err.Error())
	}

	if err = os.Rename(tmpBin, bin); err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	if _, err = os.Stat(bin); err != nil {
		return jsvm.Callback(errCb, err.Error())
	}

	go func() {
		if err := exect(os.Stdout, bin, binArgs); err != nil {
			jsvm.Callback(errCb, err)
		}

	}()

	actionTime["reload"] = time.Now()
	return otto.UndefinedValue()
}

func exect(stdout io.Writer, icmdl ...interface{}) error {
	cmdl := formatArgs(icmdl...)
	c := exec.Command(cmdl[0], cmdl[1:]...)
	stderr := &bytes.Buffer{}
	c.Stderr = stderr
	if stdout != nil {
		c.Stdout = stdout
	}

	c.Run()
	if stderr.Len() != 0 {
		return errors.New(stderr.String())
	}
	return nil
}

func formatArgs(iargs ...interface{}) (args []string) {

	for _, iarg := range iargs {
		switch arg := iarg.(type) {
		case string:
			args = append(args, arg)
		case []string:
			args = append(args, arg...)
		case otto.Value:
			iiarg, _ := arg.Export()
			switch arg := iiarg.(type) {
			case string:
				args = append(args, arg)
			case []string:
				args = append(args, arg...)
			}
		}
	}
	return
}

func getBinName(path_, name string) string {
	binName := path.Join(path_, name)
	switch runtime.GOOS {
	case "windows":
		binName = binName + ".exe"

	}
	return binName
}
