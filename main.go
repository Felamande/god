package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var version = "1.3.1"

func main() {
	app := cli.NewApp()
	app.Version = version
	app.Name = filepath.Base(os.Args[0])
	app.HelpName = app.Name
	app.Usage = "An automation tool for Go project."
	app.Commands = append(app.Commands, appRun, initCmd)
	app.Run(os.Args)
}
