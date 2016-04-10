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

type watchTaskRunner struct {
	wildcard string
	eventCb  otto.Value
	// errCb    otto.Value
	lastTime time.Time
	lastPath string
	unique   bool
}

var allTasks map[string]*watchTaskRunner
var SubCmd map[string]jsvm.Func
var ignored []string
var wd string
var Init jsvm.Func

var onceInitMod = new(sync.Once)

func init() {
	onceInitMod.Do(func() {
		allTasks = make(map[string]*watchTaskRunner)
		SubCmd = make(map[string]jsvm.Func)
		wd, _ = os.Getwd()
		wd = filepath.ToSlash(wd)
		if p := jsvm.Module("god"); p != nil {
			p.Extend("watch", watch)
			p.Extend("ignore", ignore)
			p.Extend("init", initfn)
			p.Extend("subcmd", subcmd)
		}
	})

}

func GetTask(name string) *watchTaskRunner {
	return allTasks[name]
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
			ignored = append(ignored, arg)
		case []string:
			ignored = append(ignored, arg...)
		}
	}
	return otto.UndefinedValue()
}

func watch(call otto.FunctionCall) otto.Value {
	name := call.Argument(0).String()
	wildcard := call.Argument(1).String()
	unique, _ := call.Argument(2).ToBoolean()
	eventCb := call.Argument(3)
	allTasks[name] = &watchTaskRunner{wildcard, eventCb, time.Now(), "", unique}

	return otto.UndefinedValue()
}

func subcmd(call otto.FunctionCall) otto.Value {
	cmdName := call.Argument(0).String()
	CbFunc := call.Argument(1)
	SubCmd[cmdName] = jsvm.Func(CbFunc)
	return otto.UndefinedValue()
}

var walkOnce = new(sync.Once)
var w *fsnotify.Watcher

func BeginWatch(taskNames ...string) {
	tasks := make(map[string]*watchTaskRunner)

	if len(taskNames) == 1 && taskNames[0] == "*" {
		tasks = allTasks
	} else {
		for _, tn := range taskNames {
			if t, ok := allTasks[tn]; ok {
				tasks[tn] = t
			}

		}
	}

	if len(tasks) == 0 {
		return
	}
	walkOnce.Do(func() {
		w, _ = fsnotify.NewWatcher()
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			for _, ig := range ignored {
				// fmt.Println(filepath.ToSlash(path), ig, "watched",wildmatch.IsSubsetOf(filepath.ToSlash(path), ig))
				if wildmatch.IsSubsetOf(filepath.ToSlash(path), ig) {
					// fmt.Println(filepath.ToSlash(path), ig, "skipped",wildmatch.IsSubsetOf(filepath.ToSlash(path), ig))

					return filepath.SkipDir
				}
			}
			if !info.IsDir() {
				return nil
			}
			w.Add(path)

			return nil
		})
	})

	for {
		select {

		case e := <-w.Events:
			Name := filepath.ToSlash(e.Name)
			rel, abs := getPath(Name)
			dir := getDir(rel)

			switch e.Op {
			case fsnotify.Create:
				if isDir(abs) {
					w.Add(rel)
				}
			case fsnotify.Write:
				var uniqueTask *watchTaskRunner
				var normalTasks []*watchTaskRunner
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
					}

				}
				if uniqueTask != nil {
					uniqueTask.raise(abs, rel, dir)
					for _, nt := range normalTasks {
						nt.delay(rel)
					}
					continue
				}
				for _, t := range normalTasks {
					t.raise(abs, rel, dir)
				}
			default:

			}

		}
	}
	// return otto.UndefinedValue()

}

func (t *watchTaskRunner) raise(abs, rel, dir string) {
	jsvm.Callback(t.eventCb, jsvm.ToObject(jsvm.O{"rel": rel, "abs": abs, "dir": dir}))
	t.delay(rel)
}

func (t *watchTaskRunner) delay(rel string) {
	t.lastPath = rel
	t.lastTime = time.Now()
}

func (t *watchTaskRunner) match(rel string) bool {
	p := path.Clean(rel)
	if strings.HasSuffix(p, "..") {
		return false
	}
	p = strings.TrimLeft(p, "./")
	return wildmatch.IsSubsetOf(p, t.wildcard)
}

func (t *watchTaskRunner) intervalTooShort(rel string) bool {
	// fmt.Println(rel, time.Now().Sub(t.lastTime).Seconds())
	return t.lastPath == rel && time.Now().Sub(t.lastTime).Seconds() < 2
}

func getPath(raw string) (rel, abs string) {
	if !filepath.IsAbs(raw) {
		rel = formatRel(raw)
		abs = path.Join(wd, rel)
		return
	}

	tmp := strings.Split(abs, wd)
	if len(tmp) == 1 {
		rel = "."
	} else {
		rel = formatRel(tmp[1])
	}
	return
}

func formatRel(path string) string {
	if path != "." && path != ".." && !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		path = "./" + path
	}
	return path
}
func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func getDir(p string) string {
	if filepath.IsAbs(p) {
		return path.Dir(p)
	}
	var suffix string
	dir := path.Dir(p)
	if dir != "." && dir != ".." {
		suffix = "./"
	}
	return suffix + dir
}
