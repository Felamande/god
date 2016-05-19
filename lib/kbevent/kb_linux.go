
package kbevent

import (
	// "gopkg.in/xkg.v0"
)

func New() Interface {
	return &Keyboard{}
}

type Keyboard struct {
}

func (k *Keyboard) Init() (err error ){
	return nil
}

func (k *Keyboard) KeyCodeOf(key string) (uint8, bool) {
	
	return 0, false
}

func (k *Keyboard) ModifierCodeOf(m string) (uint8, bool) {
	return 0, false
}

func (k *Keyboard) Parse(seq string) (mods []uint8, key uint8, err error) {

	return nil, 0, ErrUnimplemeted
}

func (k *Keyboard) GetSeq(mods []uint8, keyCode uint8) (string, error) {
	return "", ErrUnimplemeted

}

func (k *Keyboard) Call(seq string) error {
	return ErrUnimplemeted

}

func (k *Keyboard) HandlerOf(seq string) func() {
	return nil

}

func (k *Keyboard) Bind(seq string, f func()) error {
	return ErrUnimplemeted

}

func (k *Keyboard) ReadEvents(seqChan chan string, errChan chan error) {

	errChan <- ErrUnimplemeted
}

func (k *Keyboard) Start(errChan chan error) {

	errChan <- ErrUnimplemeted

}
