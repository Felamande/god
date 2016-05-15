package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/codegangsta/cli"
	"gopkg.in/fsnotify.v1"

	"github.com/Felamande/god/modules/god"

	"github.com/Felamande/god/lib/pathutil"
)

// "github.com/Felamande/god/lib/kbevent"

func watchCmd() cli.Command {
	return cli.Command{
		Name:   "watch",
		Action: cmder.watch,
		Usage:  "watch taskName...(use * to watch all)",
	}
}

func unwatchCmd() cli.Command {
	return cli.Command{
		Name:   "unwatch",
		Action: cmder.unwatch,
		Usage:  "unwatch tasks, or use ctrl+alt+p",
	}
}

func (c *Cmder) watch(ctx *cli.Context) error {

	go c.BeginWatch(ctx.Args()...)
	return nil

	// go c.BeginWatch(ctx.Args()...)
}

func (c *Cmder) unwatch(ctx *cli.Context) error {
	if c.isStartWatch {
		go func() {
			c.stopWChan <- true
		}()
	}
	return nil
}

func (c *Cmder) started() {
	c.isStartWatch = true
}

func (c *Cmder) stopped() {
	c.isStartWatch = false
}

func (c *Cmder) BeginWatch(taskNames ...string) {
	c.started()
	defer c.stopped()

	pathutil.SetPrefix(pathutil.PrefixDotSlash)
	tasks := []*god.WatchTaskRunner{}

	if len(taskNames) == 1 && taskNames[0] == "*" {
		tasks = god.GetAllTasks()
	} else {
		for _, tn := range taskNames {
			if t := god.GetTask(tn); t != nil {
				tasks = append(tasks, t)
			}

		}
	}

	if len(tasks) == 0 {
		return
	}

	c.w, _ = fsnotify.NewWatcher()
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if god.IsIgnore(filepath.ToSlash(path)) {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			return nil
		}
		c.w.Add(path)

		return nil
	})
loop:
	for {
		select {
		case <-c.stopWChan:
			break loop

		case e := <-c.w.Events:
			abs, rel, _ := pathutil.GetPath(e.Name)
			dir := "./" + path.Dir(rel)

			switch e.Op {
			case fsnotify.Create:
				if isDir(abs) {
					c.w.Add(rel)
				}
			case fsnotify.Write:
				var uniqueTask *god.WatchTaskRunner
				var normalTasks []*god.WatchTaskRunner
				for _, t := range tasks {

					if t.Match(rel) {
						if t.IntervalTooShort(rel) {
							t.Delay(rel)
							continue
						}
						if t.Unique() {
							uniqueTask = t
						} else {
							normalTasks = append(normalTasks, t)
						}
					}

				}
				if uniqueTask != nil {
					uniqueTask.Raise(abs, rel, dir)
					for _, nt := range normalTasks {
						nt.Delay(rel)
					}
					continue
				}
				for _, t := range normalTasks {
					t.Raise(abs, rel, dir)
				}
			default:

			}

		}
	}

	err := c.w.Close()
	fmt.Printf("stop watch %v", err)
	// return otto.UndefinedValue()

}

func (c *Cmder) StopWatch() {
	c.stopWChan <- true
}

func isDir(p string) bool {
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}
