//go:build !windows

package autostart

import "errors"

func Available() bool { return false }

func Status() (bool, error) {
	return false, errors.New("Windows startup is only available on Windows")
}

func Set(bool) (bool, error) {
	return false, errors.New("Windows startup is only available on Windows")
}
