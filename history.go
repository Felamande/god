package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var historyCmd = cli.Command{
	Name:  "history",
	Usage: "show history",
	Subcommands: []cli.Command{
		{
			Name:   "clean",
			Usage:  "clean hisory.",
			Action: cmder.cleanHistory,
		},
	},
}

func (c *Cmder) cleanHistory(ctx *cli.Context) {
	if c.history == nil {

		return
	}
	err := c.history.Truncate(0)
	if err != nil {
		fmt.Println(err)
	}
	c.history.Sync()

}
