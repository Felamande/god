package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"

	"github.com/Felamande/god/lib/jsvm"

	"github.com/codegangsta/cli"
	"github.com/peterh/liner"

	"github.com/Felamande/god/lib/pathutil"
	_ "github.com/Felamande/god/modules/go"
	"github.com/Felamande/god/modules/god"
	"github.com/Felamande/god/modules/hotkey"
	_ "github.com/Felamande/god/modules/localstorage"

	_ "github.com/Felamande/god/modules/log"
	_ "github.com/Felamande/god/modules/os"
	_ "github.com/Felamande/god/modules/path"
)

var version = "2.0.0-beta3"

type vars map[string]string

//Cmder command manager
type Cmder struct {
	line           *liner.State
	history        *os.File
	w              *fsnotify.Watcher
	stopWChan      chan bool
	watchm         *sync.Mutex
	wd             string
	isStartWatch   bool
	global         vars
	varsReplacer   *strings.Replacer
	escapeReplacer *strings.Replacer
}

var cmder *Cmder

func init() {

	cmder = &Cmder{
		stopWChan: make(chan bool, 1),
		watchm:    new(sync.Mutex),
		global:    vars{"ps1": "(god) "},
		varsReplacer: strings.NewReplacer(
			"${wd}", pathutil.Wd(),
		),
		escapeReplacer: strings.NewReplacer(
			`\s`, " ",
			`\n`, "\n",
		),
	}

}
func main() {
	err := jsvm.Run("god.js")
	if err != nil {
		fmt.Println(err)
	}

	app := cli.NewApp()
	app.Version = version
	app.Name = filepath.Base(os.Args[0])
	app.HelpName = app.Name
	app.Usage = "An automation tool for Go project."
	app.Action = cmder.enterAction
	app.Commands = append(app.Commands, initCmd(), watchCmd(),
		reloadCmd(), exitCmd(), historyCmd(), unwatchCmd(), setCmd())

	for name, CbFunc := range god.SubCmd {
		app.Commands = append(app.Commands, newSubCmd(name, CbFunc))
	}

	err = hotkey.Bind("ctrl+d", func() { cmder.exit(nil) })
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

		argStr = c.evalGlobal(argStr)
		splits := strings.Split(argStr, " ")
		var args = append([]string{}, "god")
		for _, s := range splits {
			if s != "" {
				args = append(args, s)
			}
		}
		if len(args) == 1 {
			continue
		}
		if ctx.App.Command(args[1]) == nil {
			continue
		}

		ctx.App.Run(args)

	}

}

func (c *Cmder) evalGlobal(raw string) string {
	return c.varsReplacer.Replace(raw)
}

func (c *Cmder) evalEscape(raw string) string {
	return c.escapeReplacer.Replace(raw)
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

func (v vars) String() string {
	if len(v) == 0 {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)

}

func setCmd() cli.Command {
	return cli.Command{
		Name:      "set",
		Usage:     fmt.Sprintf("set <key> <value> defaults:%v\n", cmder.global),
		Action:    cmder.setfn,
		UsageText: `use \s to represent whitespace`,
	}
}

func (c *Cmder) setfn(ctx *cli.Context) {
	if len(ctx.Args()) != 2 {
		fmt.Println("set key value")
		return
	}
	k := ctx.Args()[0]
	v := ctx.Args()[1]
	c.global[k] = c.evalEscape(v)
}
