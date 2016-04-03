package god

import (
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
			if e.Op == fsnotify.Rename || e.Op == fsnotify.Remove || e.Op == fsnotify.Chmod {
				continue
			}
			rel := format(e.Name)
			abs := filepath.Join(wd, rel)
			if isDir(abs) {
				w.Add(rel)
			}
			var uniqueTask *taskRunner
			var normalTasks []*taskRunner
			for _, t := range tasks {

				if t.match(rel) {
					if t.intervalTooShort(rel) {
						t.delay(rel)
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
				// fmt.Println("I am unique")
				uniqueTask.raise(abs, rel)
				for _, nt := range normalTasks {
					nt.delay(rel)
				}
				continue
			}
			for _, t := range normalTasks {
				// fmt.Println("I am not unique")
				t.raise(abs, rel)
			}
		}
	}
	// return otto.UndefinedValue()

}

func (t *taskRunner) raise(abs, rel string) {
	jsvm.Callback(t.eventCb, abs, rel)
	t.delay(rel)
}

func (t *taskRunner) delay(rel string) {
	t.lastPath = rel
	t.lastTime = time.Now()
}

func (t *taskRunner) match(rel string) bool {
	return wildmatch.IsSubsetOf(path.Clean(rel), t.wildcard)
}

func (t *taskRunner) intervalTooShort(rel string) bool {
	return t.lastPath == rel && time.Now().Sub(t.lastTime).Seconds() < 1
}

func format(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}
func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
