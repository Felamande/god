package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "An automation tool for Go project."
	app.Commands = append(app.Commands, appRun)
	app.Run(os.Args)
}
