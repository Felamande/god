package kbevent

import (
	"fmt"
	"runtime"
)

type Interface interface {
	Parse(seq string) (mods []uint8, key uint8, err error)

	GetSeq(mods []uint8, key uint8) (seq string, err error)

	Init() error

	KeyCodeOf(k string) (code uint8, exist bool)

	ModifierCodeOf(m string) (code uint8, exist bool)

	Call(seq string) error

	HandlerOf(seq string) func()

	Bind(seq string, f func()) error

	Start() error

	ReadEvents(chan string, chan error) error
}

var ErrUnimplemeted = fmt.Errorf("unimplemented on %s!", runtime.GOOS)
