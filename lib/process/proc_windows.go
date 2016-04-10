package process

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func killByName(name string) error {
	if !strings.HasSuffix(name, ".exe") {
		name = name + ".exe"
	}
	c := exec.Command("taskkill", "/im", name, "/f")
	stderr := bytes.NewBuffer([]byte{})
	c.Stderr = stderr
	c.Run()
	if stderr.Len() != 0 {
		return errors.New(strings.TrimSuffix(stderr.String(), "\n"))
	}
	return nil
}
