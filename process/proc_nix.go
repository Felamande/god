// +build !windows

package process

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func killByName(name string) error {

	c := exec.Command("pkill", "-9", name)
	stderr := bytes.NewBuffer([]byte{})
	c.Stderr = stderr
	c.Run()
	if stderr.Len() != 0 {
		return errors.New(strings.TrimSuffix(stderr.String(), "\n"))
	}
	return nil
}
