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
	},
}

func initfn(c *cli.Context) {

	if fi, _ := os.Stat("god.js"); fi != nil {
		if c.Bool("override") {
			ioutil.WriteFile("god.js", []byte(defaultjs), 0777)
			fmt.Println("override god.js")
		} else {
			fmt.Println("You've already had a god.js.")
		}
	} else {
		ioutil.WriteFile("god.js", []byte(defaultjs), 0777)
		fmt.Println("create god.js")
	}

}

const defaultjs = `
god = require("god")
log = require("log")
path = require("path")
os = require("os")
go  = require("go")

god.ignore(".git", ".vscode")

god.init(function(){go.reload(".")})

god.watch(["*_test.go", "**/*_test.go"],true,
    function(abs, rel) {
            go.test(path.dir(rel),[],function(err){log.error(err)})    
    }
)

god.watch("**/*.go",false,
    function(abs, rel) {
        go.install(path.dir(rel),[],function(err){log.error(err)})
    }
)

god.watch(["*.go"],false,
    function(abs, rel) {
        log.info("reload",rel)  
        go.reload(".",[],[],function(err){if(err){log.error(err)}})
    }
)
`
