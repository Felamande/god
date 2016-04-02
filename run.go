package main

import (
	"time"

	_ "github.com/Felamande/god/modules/go"
	"github.com/Felamande/god/modules/god"
	_ "github.com/Felamande/god/modules/log"
	"github.com/Felamande/jsvm"
	_ "github.com/Felamande/jsvm/module/os"
	_ "github.com/Felamande/jsvm/module/path"
	"github.com/codegangsta/cli"
	"gopkg.in/fsnotify.v1"
)

var lastReloadTime time.Time

var appRun = cli.Command{
	Name:   "run",
	Usage:  "run god configure file, default god.js in work directory.",
	Action: run,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "god.js",
			Usage: "config file.",
		},
	},
}

func run(c *cli.Context) {
	f := c.String("config")
	jsvm.Run(f)
	lastReloadTime = time.Now()
	go MustWatch(f)
	god.Run()

}

func MustWatch(file string) {
	if err := watch(file); err != nil {
		panic(err)
	}
}

func Watch(file string) error {
	return watch(file)
}

func watch(file string) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = w.Add(file)
	if err != nil {
		return err
	}
	for {
		select {
		case e := <-w.Events:
			if e.Op == fsnotify.Remove || e.Op == fsnotify.Rename || e.Op == fsnotify.Chmod {
				continue
			}
			if time.Now().Sub(lastReloadTime) < 3 {
				continue
			}
			reload(file)
			lastReloadTime = time.Now()
		}
	}
}

func reload(file string) error {
	err := jsvm.Run(file)
	if err != nil {
		return err
	}
	return nil
}
