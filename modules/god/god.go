package god

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Felamande/god/lib/jsvm"
	"github.com/Felamande/otto"
	"github.com/demon-xxi/wildmatch"
)

type WatchTaskRunner struct {
	wildcard string
	eventCb  otto.Value
	// errCb    otto.Value
	lastTime time.Time
	lastPath string
	unique   bool
}

var allTasks map[string]*WatchTaskRunner
var SubCmd map[string]jsvm.Func
var ignored map[string]bool
var wd string
var Init jsvm.Func

func init() {
	allTasks = make(map[string]*WatchTaskRunner)
	SubCmd = make(map[string]jsvm.Func)
	ignored = make(map[string]bool)
	wd, _ = os.Getwd()
	wd = filepath.ToSlash(wd)
	if p := jsvm.Module("god"); p != nil {
		p.Extend("watch", watch)
		p.Extend("ignore", ignore)
		p.Extend("init", initfn)
		p.Extend("subcmd", subcmd)
	}

}

func GetTask(name string) *WatchTaskRunner {

	return allTasks[name]
}

func GetAllTasks() (all []*WatchTaskRunner) {
	for _, t := range allTasks {
		all = append(all, t)
	}
	return
}

// function init(initfn, fnArgs...)
func initfn(call otto.FunctionCall) otto.Value {
	Init = jsvm.Func(call.Argument(0))
	return otto.UndefinedValue()
}

func ignore(call otto.FunctionCall) otto.Value {
	for _, v := range call.ArgumentList {
		iarg, _ := v.Export()
		switch arg := iarg.(type) {
		case string:
			ignored[arg] = true
		case []string:
			for _, a := range arg {
				ignored[a] = true
			}

		}
	}
	return otto.UndefinedValue()
}

func watch(call otto.FunctionCall) otto.Value {
	name := call.Argument(0).String()
	wildcard := call.Argument(1).String()
	unique, _ := call.Argument(2).ToBoolean()
	eventCb := call.Argument(3)
	allTasks[name] = &WatchTaskRunner{wildcard, eventCb, time.Now(), "", unique}

	return otto.UndefinedValue()
}

func subcmd(call otto.FunctionCall) otto.Value {
	cmdName := call.Argument(0).String()
	CbFunc := call.Argument(1)
	SubCmd[cmdName] = jsvm.Func(CbFunc)
	return otto.UndefinedValue()
}

func IsIgnore(path string) bool {
	return ignored[path]
}

func (t *WatchTaskRunner) Unique() bool {
	return t.unique
}

func (t *WatchTaskRunner) Raise(abs, rel, dir string) {
	jsvm.Callback(t.eventCb, jsvm.ToObject(jsvm.O{"rel": rel, "abs": abs, "dir": dir}))
	t.Delay(rel)
}

func (t *WatchTaskRunner) Delay(rel string) {
	t.lastPath = rel
	t.lastTime = time.Now()
}

func (t *WatchTaskRunner) Match(rel string) bool {
	p := path.Clean(rel)
	if strings.HasSuffix(p, "..") {
		return false
	}
	p = strings.TrimLeft(p, "./")
	return wildmatch.IsSubsetOf(p, t.wildcard)
}

func (t *WatchTaskRunner) IntervalTooShort(rel string) bool {
	// fmt.Println(rel, time.Now().Sub(t.lastTime).Seconds())
	return t.lastPath == rel && time.Now().Sub(t.lastTime).Seconds() < 2
}
