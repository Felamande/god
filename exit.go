package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var exitCmd = cli.Command{
	Name:   "exit",
	Action: cmder.exit,
	Usage:  "exit god, or use ctrl+D.",
}

func (c *Cmder) exitOn(*cli.Context) {

}
func (c *Cmder) exit(*cli.Context) {
	if c.line != nil && c.history != nil {
		c.line.WriteHistory(c.history)
	}

	os.Exit(0)
}
