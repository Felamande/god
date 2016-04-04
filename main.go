package main

import (
	"fmt"
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

func execFlag(c *cli.Context, name string, fn func(value interface{}) error) error {
	v, _ := lookup(c, name)
	return fn(v)
}
func isSetFlag(c *cli.Context, name string, isSetfn func(value interface{}) error, notSetfn func() error) error {
	v, isset := lookup(c, name)
	if isset {
		return isSetfn(v)
	}
	return notSetfn()

}

func noFlag(c *cli.Context, fn func() error) (bool, error) {
	if c.NumFlags() == 0 {
		return true, fn()
	}
	return false, nil
}

func lookup(c *cli.Context, name string) (value interface{}, isset bool) {
	return c.Generic(name), c.IsSet(name)
}

type flagError struct {
	flag string
	msg  string
}

func (e *flagError) Error() string {
	return fmt.Sprintf("%s: %s", e.flag, e.msg)
}
func (e *flagError) String() string {
	return fmt.Sprintf("%s: %s", e.flag, e.msg)
}
