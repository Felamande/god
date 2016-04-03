package god

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Felamande/jsvm"
	"github.com/Felamande/otto"
	"github.com/demon-xxi/wildmatch"
	"gopkg.in/fsnotify.v1"
)

type taskRunner struct {
	wildcard string
	eventCb  otto.Value
	// errCb    otto.Value
	lastTime time.Time
	lastPath string
	unique   bool
}

var tasks []*taskRunner
var ignored []string
var wd string
var initCall []otto.Value
var once = new(sync.Once)

func callInit() {
	once.Do(func() {
		lenArg := len(initCall)
		if initCall != nil && lenArg > 0 {
			iarg := make([]interface{}, len(initCall)-1)

			for idx, arg := range initCall {
				if idx == 0 {
					continue
				}
				iarg[idx-1] = arg
			}
			jsvm.Callback(initCall[0], iarg...)
		}
	})

}

func init() {
	wd, _ = os.Getwd()
	wd = format(wd)
	if p := jsvm.Module("god"); p != nil {
		p.Extend("watch", watch)
		p.Extend("ignore", ignore)
		p.Extend("init", initfn)
	}
}

// function init(initfn, fnArgs...)
func initfn(call otto.FunctionCall) otto.Value {
	initCall = call.ArgumentList
	return otto.UndefinedValue()
}

func ignore(call otto.FunctionCall) otto.Value {
	for _, v := range call.ArgumentList {
		iarg, _ := v.Export()
		switch arg := iarg.(type) {
		case string:
			ignored = append(ignored, arg)
		case []string:
			ignored = append(ignored, arg...)
		}
	}
	return otto.UndefinedValue()
}

func watch(call otto.FunctionCall) otto.Value {
	wildcards, _ := call.Argument(0).Export()
	unique, _ := call.Argument(1).ToBoolean()
	eventCb := call.Argument(2)
	switch w := wildcards.(type) {
	case string:
		tasks = append(tasks, &taskRunner{w, eventCb, time.Now(), "", unique})
	case []string:
		for _, wildcard := range w {
			tasks = append(tasks, &taskRunner{wildcard, eventCb, time.Now(), "", unique})
		}
	}

	return otto.UndefinedValue()
}

func Run() {
	callInit()

	w, _ := fsnotify.NewWatcher()
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		for _, ig := range ignored {
			if wildmatch.IsSubsetOf(format(path), ig) {

				return filepath.SkipDir
			}
		}
		if !info.IsDir() {
			return nil
		}
		w.Add(path)

		return nil
	})
	for {
		select {
		case e := <-w.Events:
			rel := format(e.Name)
			abs := filepath.Join(wd, rel)
			if isDir(abs) {
				w.Add(rel)
			}
			var uniqueTask *taskRunner
			var normalTasks []*taskRunner
			for _, t := range tasks {

				if wildmatch.IsSubsetOf(path.Clean(rel), t.wildcard) {
					if t.lastPath == rel && time.Now().Sub(t.lastTime).Seconds() < 1 {
						t.lastPath = rel
						t.lastTime = time.Now()
						continue
					}
					if t.unique {
						uniqueTask = t
					} else {
						normalTasks = append(normalTasks, t)
					}
					// jsvm.Callback(t.eventCb, abs, rel)
				}

			}
			if uniqueTask != nil {
				fmt.Println("I am unique")
				jsvm.Callback(uniqueTask.eventCb, abs, rel)
				uniqueTask.lastPath = rel
				uniqueTask.lastTime = time.Now()
				continue
			}
			for _, t := range normalTasks {
				fmt.Println("I am not unique")
				jsvm.Callback(t.eventCb, abs, rel)
				t.lastPath = rel
				t.lastTime = time.Now()
			}
		}
	}
	// return otto.UndefinedValue()

}
func format(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
