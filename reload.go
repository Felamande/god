package main

import (
	"github.com/Felamande/god/lib/jsvm"
	"github.com/Felamande/god/modules/god"
	"github.com/codegangsta/cli"
)

func reloadCmd() cli.Command {
	cli.Command{
		Name:   "reload",
		Usage:  "reload god.js",
		Action: reload,
	}
}

func reload(c *cli.Context) {
	jsvm.Run("god.js")
	god.Init.Call()
}
