package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/codegangsta/cli"
)

func execFlag(c *cli.Context, name string, fn func(value flag.Value) error) error {
	v, _ := lookup(c, name)
	return fn(v)
}

func isSetFlag(c *cli.Context, name string, isSetfn func(value flag.Value) error) error {
	v, isset := lookup(c, name)
	if isset && v != nil {
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

func mustHaveCond(c *cli.Context, name string, flagIsSet bool, cond []string, execf func(flag.Value) error, miss func() error) (bool, error) {
	v, isSet := lookup(c, name)
	if flagIsSet && !isSet {
		return false, &flagError{name, "flag is not set."}
	}

	for _, f := range cond {
		if !c.IsSet(f) {
			return false, miss()
		}
	}
	return true, execf(v)
}

func lookup(c *cli.Context, name string) (value flag.Value, isset bool) {

	return c.Generic(name).(flag.Value), c.IsSet(name)
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

func parseArgs(args []string, prefix string) (flags map[string]interface{}, nargs []string, err error) {
	flags = make(map[string]interface{})

	for _, arg := range args {
		alen := len(arg)
		switch {
		case arg[0] == '"':
			fallthrough
		case arg[0] == '\'':
			if arg[alen-1] != arg[0] {
				return nil, nil, fmt.Errorf(`syntax error: expect %s`, string(arg[0]))
			}
			nargs = append(nargs, arg[1:alen-1])
		case strings.HasPrefix(arg, prefix):
			expr := strings.TrimLeft(arg, prefix)

			k, v, err := eval(expr)
			if err != nil {
				return nil, nil, err
			}

			flags[k] = v
		default:
			nargs = append(nargs, arg)
		}
	}
	return
}

func parseRaw(cmdline string, prefix string) (flags map[string]interface{}, nargs []string, err error) {

	return
}

func eval(expr string) (k string, v interface{}, err error) {
	length := len(expr)
	if expr[length-1] == '=' {
		return "", "", errors.New("syntax error: expect value")
	}
	var i int
	for i = range expr {
		if expr[i] == '=' {
			break
		}
	}
	// fmt.Println("i =", i)
	if i == length-1 {
		k = expr
		v = true
	} else {
		k, v = expr[:i], expr[i+1:]
	}
	return k, v, nil
}
