package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

var initCmd = cli.Command{
	Name:   "init",
	Usage:  "create a default god.js",
	Action: initfn,
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

func initfn(c *cli.Context) {

	noflag, err := noFlag(c, func() error {
		if fi, _ := os.Stat("god.js"); fi != nil {
			fmt.Println("already had a god.js, use --override or -o to override it with the default one.")
			return nil
		}
		return ioutil.WriteFile("god.js", []byte(defaultjs), 0777)

	})

	if noflag {
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	isSetFlag(c, "override",
		func(_ interface{}) error {
			ioutil.WriteFile("god.js", []byte(defaultjs), 0777)
			fmt.Println("--override: god.js overrided.")
			return nil
		},
	)

	isSetFlag(c, "ignore",
		func(_ interface{}) error {
			f, err := os.OpenFile(".gitignore", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
			if err != nil {
				fmt.Println("--ignore:", err)
				return nil
			}
			defer f.Close()

			f.Seek(0, SEEK_END)
			f.Seek(-1, SEEK_CUR)

			buf := make([]byte, 1)
			f.Read(buf)

			newline := ""
			if buf[0] != '\n' {
				newline = "\n"
			}

			_, err = f.WriteString(newline + "god.js\n")
			if err != nil {
				fmt.Println("--ignore:", err)
				return nil
			}
			fmt.Println("--ignore: ok")
			return nil
		},
	)

}

const SEEK_END = 2
const SEEK_CUR = 1

const defaultjs = `

//modules located in github.com/Felamande/jsvm/modules 
//and github.com/Felamande/god/modules
//you can write modules yourself if you're familiar witch otto.
god = require("god")
log = require("log")
path = require("path")
os = require("os")
go = require("go")

//changes of ignored files or dirs will not be watched 
god.ignore(".git", ".vscode")

var buildArgs = []
var installArgs = []
var binArgs = []
var testArgs = [] //args for reloaded binaries 

//will be call after god starts and before god watches changes.
god.init(function() { go.reload(".", buildArgs, binArgs)})

// event 
// event.rel, relative path of matched file or directory
// event.abs, absolute path of matched file or directory
// event.dir, relative parent directory of matched file or directory
// path seperator will be slash on windows.
god.watch(["*_test.go", "**/*_test.go"], true,
    function(event) {
        go.test(event.dir, testArgs, function(err) { log.error(err) })
    }
)

// ** will match ONE or more directories
// * will match just ONE directory or as many chars as possible except slash .
god.watch("**/*.go", false,
    function(event) {
        go.install(event.dir, installArgs, function(err) { log.error(err) })
    }
)

god.watch(["*.go"], false,
    function(event) {
        log.info("reload", event.rel, )
        go.reload(".", buildArgs, binArgs, function(err) { if (err) { log.error(err) } })
    }
)
`
