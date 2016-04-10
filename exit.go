package main

import (
	"os"

	"github.com/codegangsta/cli"
)

var exitCmd = cli.Command{
	Name:   "exit",
	Action: exit,
	Usage:  "exit god, or use ctrl+D.",
}

func exit(*cli.Context) {
	os.Exit(0)
}
