package main

import (
	"errors"

	"github.com/codegangsta/cli"
)

func historyCmd() cli.Command {
	return cli.Command{
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
}

func (c *Cmder) cleanHistory(ctx *cli.Context) error {
	if c.history == nil {
		return errors.New("no history found")
	}
	err := c.history.Truncate(0)
	if err != nil {
		return err
	}
	return c.history.Sync()

}
