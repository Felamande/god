package main

import (
	"github.com/Felamande/god/modules/god"
	"github.com/Felamande/jsvm"
	"github.com/codegangsta/cli"
)

var reloadCmd = cli.Command{
	Name:   "reload",
	Usage:  "reload god.js",
	Action: reload,
}

func reload(c *cli.Context) {
	jsvm.Run("god.js")
	god.Init.Call()
}
