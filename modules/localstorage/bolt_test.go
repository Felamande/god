package localstorage

import (
	"testing"
)

func TestGet(t *testing.T) {
	b, _ := localStorageType["bolt"].Get([]byte("greetings"))
	t.Log(string(b))
	b, _ = localStorageType["bolt"].Get([]byte("greetings"))
	t.Log(string(b))
}
