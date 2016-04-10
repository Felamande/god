package main

import (

	// "github.com/Felamande/god/lib/kbevent"

	"github.com/Felamande/god/modules/god"
	"github.com/codegangsta/cli"
)

var watchCmd = cli.Command{
	Name:   "watch",
	Action: watch,
	Usage:  "watch taskName...(use * to watch all)",
}

func watch(c *cli.Context) {
	// kb.BindStr("ctrl+shift+d", func() { os.Exit(0) })

	go god.BeginWatch(c.Args()...)
}
