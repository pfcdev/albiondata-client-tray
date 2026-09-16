//go:build !windows

package consoleviewer

import (
	"errors"
	"io"
)

func Run(string, io.Reader) error {
	return errors.New("console viewer is only available on Windows")
}
