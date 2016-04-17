package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"

	"github.com/Felamande/god/lib/jsvm"

	"github.com/codegangsta/cli"
	"github.com/peterh/liner"

	_ "github.com/Felamande/god/modules/go"
	"github.com/Felamande/god/modules/god"
	"github.com/Felamande/god/modules/hotkey"
	_ "github.com/Felamande/god/modules/log"
	_ "github.com/Felamande/god/modules/os"
	_ "github.com/Felamande/god/modules/path"
)

var version = "2.0.0-beta2"

//Cmder command manager
type Cmder struct {
	line      *liner.State
	history   *os.File
	w         *fsnotify.Watcher
	stopWChan chan bool
	watchm    *sync.Mutex
	wd        string

	global map[string]string
}

var cmder = &Cmder{
	stopWChan: make(chan bool, 1),
	watchm:    new(sync.Mutex),
	global:    make(map[string]string),
}

func init() {
	err := jsvm.Run("god.js")
	if err != nil {
		fmt.Println(err)
	}
	cmder.global["ps1"] = "(god) "
}
func main() {

	app := cli.NewApp()
	app.Version = version
	app.Name = filepath.Base(os.Args[0])
	app.HelpName = app.Name
	app.Usage = "An automation tool for Go project."
	app.Action = cmder.enterAction
	app.Commands = append(app.Commands, initCmd, watchCmd, reloadCmd, exitCmd, historyCmd)

	for name, CbFunc := range god.SubCmd {
		app.Commands = append(app.Commands, newSubCmd(name, CbFunc))
	}

	err := hotkey.Bind("ctrl+d", func() { cmder.exit(nil) })
	if err != nil {
		fmt.Println(err)

	}
	err = hotkey.Bind("ctrl+alt+p", func() { cmder.unwatch(nil) })
	if err != nil {
		fmt.Println(err)
	}

	go hotkey.ApplyAll()
	app.Run(os.Args)
}

func (c *Cmder) enterAction(ctx *cli.Context) {

	god.Init.Call()
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(false)
	c.line = line

	historyFile, err := os.OpenFile(filepath.Join(os.TempDir(), ".god_hostory"), os.O_CREATE|os.O_RDWR, 0777)

	if err != nil {
		fmt.Println(err)

	}
	c.history = historyFile
	line.ReadHistory(historyFile)

	for {
		argStr, err := line.Prompt(c.global["ps1"])
		argStr = strings.TrimSpace(argStr)
		if len(argStr) == 0 {
			continue
		}
		line.AppendHistory(argStr)
		if err == liner.ErrPromptAborted {
			os.Exit(0)
		} else if err != nil {
			fmt.Println(err)
			continue
		}

		splits := strings.Split(argStr, " ")
		var args = append([]string{}, "god")
		for _, s := range splits {
			if s != "" {
				args = append(args, s)
			}
		}
		if len(args) == 0 {
			continue
		}
		if ctx.App.Command(args[1]) == nil {
			continue
		}

		ctx.App.Run(args)

	}

}

func newSubCmd(name string, CbFunc jsvm.Func) cli.Command {
	return cli.Command{
		Name: name,
		Action: func(c *cli.Context) {
			arg := c.Args()
			flags, nargs, err := parseArgs(arg, "-")
			if err != nil {
				fmt.Println(name+": ", err)
				return
			}
			CbFunc.Call(nargs, flags)
		},
		Usage: "user defined command " + name,

		SkipFlagParsing: true,
	}
}
