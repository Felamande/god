package main

import (
	"github.com/Felamande/god/lib/jsvm"
	"github.com/codegangsta/cli"
)

func reloadCmd() cli.Command {
	return cli.Command{
		Name:   "reload",
		Usage:  "reload god.js",
		Action: reload,
	}
}

func reload(c *cli.Context) {
	jsvm.Run("god.js")
}
