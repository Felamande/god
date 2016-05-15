package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func exitCmd() cli.Command {
	return cli.Command{
		Name:   "exit",
		Action: cmder.exit,
		Usage:  "exit god, or use ctrl+D.",
	}
}

func (c *Cmder) exit(*cli.Context) error {
	if c.line != nil && c.history != nil {
		c.line.WriteHistory(c.history)
	}

	os.Exit(0)
	return nil
}
