package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func execFlag(c *cli.Context, name string, fn func(value interface{}) error) error {
	v, _ := lookup(c, name)
	return fn(v)
}

func isSetFlag(c *cli.Context, name string, isSetfn func(value interface{}) error) error {
	v, isset := lookup(c, name)
	if isset {
		return isSetfn(v)
	}
	return nil

}

func noFlag(c *cli.Context, fn func() error) (bool, error) {
	if c.NumFlags() == 0 {
		return true, fn()
	}
	return false, nil
}

func mustHaveCond(c *cli.Context, name string, flagIsSet bool, cond []string, execf func(interface{}) error, miss func() error) (bool, error) {
	if flagIsSet && !c.IsSet(name) {
		return false, &flagError{name, "flag is not set."}
	}

	for _, f := range cond {
		if !c.IsSet(f) {
			return false, miss()
		}
	}
	return true, execf(c.Generic(name))
}

func lookup(c *cli.Context, name string) (value interface{}, isset bool) {
	return c.Generic(name), c.IsSet(name)
}

type flagError struct {
	flag string
	msg  string
}

func (e *flagError) Error() string {
	return fmt.Sprintf("--%s: %s", e.flag, e.msg)
}
func (e *flagError) String() string {
	return fmt.Sprintf("--%s: %s", e.flag, e.msg)
}
