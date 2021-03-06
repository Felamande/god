package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

func initCmd() cli.Command {
	return cli.Command{
		Name:   "init",
		Usage:  "create a default god.js",
		Action: cmder.initfn,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "override, o",
				Usage: "override god.js with default.",
			},
			cli.BoolFlag{
				Name:  "ignore, i",
				Usage: "add god.js to .gitignore",
			},
		},
	}
}
func (c *Cmder) initfn(ctx *cli.Context) error {

	if noflag, err := noFlag(ctx, func() error {
		if fi, _ := os.Stat("god.js"); fi != nil {
			fmt.Println("already had a god.js, use --override or -o to override it with the default one.")
			return nil
		}
		return ioutil.WriteFile("god.js", []byte(defaultjs), 0777)

	}); err != nil {
		return err
	} else if noflag {
		return nil
	}

	isSetFlag(ctx, "override",
		func(flag.Value) error {
			ioutil.WriteFile("god.js", []byte(defaultjs), 0777)
			fmt.Println("--override: god.js overrided.")
			return nil
		},
	)

	return isSetFlag(ctx, "ignore",
		func(flag.Value) error {
			f, err := os.OpenFile(".gitignore", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
			if err != nil {
				return &flagError{flag: "ignore", msg: err.Error()}
			}
			defer f.Close()

			f.Seek(0, os.SEEK_END)
			f.Seek(-1, os.SEEK_CUR)

			buf := make([]byte, 1)
			f.Read(buf)

			newline := ""
			if buf[0] != '\n' {
				newline = "\n"
			}

			_, err = f.WriteString(newline + "god.js\n")
			if err != nil {
				return &flagError{flag: "ignore", msg: err.Error()}
			}

			return nil
		},
	)

}

const defaultjs = `//modules located in github.com/Felamande/god/modules
//you can write modules yourself if you're familiar with otto.
god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")
hk  = require("hotkey")


//bind hotkey as you like
hk.bind("ctrl+shift+o",function() {
    log.info("hotkey","ctrl+shift+o")
})
hk.bind("ctrl+shift+k",function() {
    log.info("hotkey","ctrl+shift+k")
})

//changes of ignored files or dirs will not be watched
god.ignore(".git", ".vscode")

var buildArgs = []
var installArgs = []
var binArgs = []
var testArgs = [] //args for reloaded binaries

//will be called immediately after god starts.
god.init(function() {console.log("hello")})


// define your subcommand, flags and arguments will be passed to the callback function.
// (god) subcmd "-willnot=-be-parsed" name=what stringvalue -key=value -testarg=-test.v -godebug=gctrace=1 -boolval
// will be parsed as
// nargs = ["-willnot=-be-parsed", "name=what", "stringvalue"],
// flags = {"key":"value", "testarg":"-test.v", "godebug":"gctrace=1", "boolvar":true}
god.subcmd("print",function(nargs,flags){
   log.info(JSON.stringify(nargs),JSON.stringify(flags))
})

god.subcmd("eval",function (nargs,flags) {
    console.log(eval(nargs[0]))
})

god.subcmd("test",function(pkgs,flags){
    for(i in pkgs){
       log.info("test",pkgs[i])
        go.test(pkgs[i], testArgs, function(err) { log.error(err) })
    }

})


god.subcmd("exec",function(nargs,flags){os.exec(nargs)})

// function watch(name, wildcard, isUnique, callback)
// if isUnique, the event which matches multiple wildcards will only be sent to the unique callback.
//
// function callback(event)
// event.rel, relative path of matched file or directory
// event.abs, absolute path of matched file or directory
// event.dir, relative parent directory of matched file or directory
//
// path seperator will be slash on windows.
// watch tasks will not start until you type the subcommand "watch [taskname...]",
// after that tasks will run in a goroutine.

god.watch("btest","*_test.go", true,
    function(event) {
        log.info("test",event.dir)
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

// ** will match ONE or more directories
// * will match just ONE directory or as many chars as possible except the slash.
god.watch("ptest", "**/*_test.go", true,
    function(event) {
        log.info("test",event.dir)
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

god.watch("pinstall","**/*.go", false,
    function(event) {
        log.info("install",event.dir)
        go.install(event.dir, installArgs, function(err) { log.error(err) })
    }
)

god.watch("breload","*.go", false,
    function(event) {
        log.info("reload", event.dir)
        go.reload(".", buildArgs, binArgs, function(err) { if (err) { log.error(err) } })
    }
)

// TODO:
// 2.separate normal tasks from watch tasks.
`
