package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Felamande/jsvm"

	"github.com/codegangsta/cli"
	"github.com/peterh/liner"

	_ "github.com/Felamande/god/modules/go"
	"github.com/Felamande/god/modules/god"
	"github.com/Felamande/god/modules/hotkey"
	_ "github.com/Felamande/god/modules/log"
	_ "github.com/Felamande/jsvm/module/os"
	_ "github.com/Felamande/jsvm/module/path"
)

var version = "2.0.0-beta"

func init() {
	err := jsvm.Run("god.js")
	if err != nil {
		fmt.Println(err)
	}

}
func main() {

	app := cli.NewApp()
	app.Version = version
	app.Name = filepath.Base(os.Args[0])
	app.HelpName = app.Name
	app.Usage = "An automation tool for Go project."
	app.Action = enterAction
	app.Commands = append(app.Commands, initCmd, watchCmd, reloadCmd, exitCmd)

	for name, CbFunc := range god.SubCmd {
		app.Commands = append(app.Commands, newSubCmd(name, CbFunc))
	}

	go hotkey.ApplyAll()
	app.Run(os.Args)
}

func enterAction(c *cli.Context) {
	god.Init.Call()
	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(false)

	historyFile, err := os.OpenFile(filepath.Join(os.TempDir(), ".god_hostory"), os.O_CREATE|os.O_RDWR, 0777)
	if err == nil {
		line.ReadHistory(historyFile)
	}

	for {
		argStr, err := line.Prompt(`(god) `)
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
		if c.App.Command(args[1]) == nil {
			continue
		}

		c.App.Run(args)

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
