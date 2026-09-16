//go:build windows

// Package autostart manages the same logon task that the NSIS installer
// creates. The task retains highest privileges for packet capture.
package autostart

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/windows"
)

const taskName = "Albion Data Client"

var setMu sync.Mutex

func Available() bool { return true }

// Status reports whether Windows will launch the client at logon.
func Status() (bool, error) {
	_, enabled, err := queryTask()
	return enabled, err
}

func queryTask() (exists, enabled bool, err error) {
	cmd := exec.Command("schtasks.exe", "/Query", "/TN", taskName, "/XML")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW}
	data, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return false, false, nil // The task has not been installed yet.
		}
		return false, false, err
	}
	enabled, err = enabledFromTaskXML(data)
	return true, enabled, err
}

// Set enables/disables the installed task, creating it when first enabled.
// A non-elevated GUI asks Windows for UAC consent to change the task.
func Set(enabled bool) (bool, error) {
	setMu.Lock()
	defer setMu.Unlock()

	exists, current, err := queryTask()
	if err != nil {
		return false, err
	}
	if current == enabled && (exists || !enabled) {
		return current, nil
	}

	args := []string{"/Change", "/TN", taskName}
	if enabled {
		if exists {
			args = append(args, "/Enable")
		} else {
			exe, err := os.Executable()
			if err != nil {
				return false, err
			}
			args = []string{"/Create", "/F", "/SC", "ONLOGON", "/RL", "HIGHEST", "/TN", taskName, "/TR", fmt.Sprintf(`"%s" -minimize`, exe)}
		}
	} else {
		args = append(args, "/Disable")
	}

	if windows.GetCurrentProcessToken().IsElevated() {
		cmd := exec.Command("schtasks.exe", args...)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: windows.CREATE_NO_WINDOW}
		if output, err := cmd.CombinedOutput(); err != nil {
			return current, fmt.Errorf("could not change Windows startup: %w: %s", err, strings.TrimSpace(string(output)))
		}
	} else {
		verb, _ := windows.UTF16PtrFromString("runas")
		file, _ := windows.UTF16PtrFromString(filepath.Join(os.Getenv("SystemRoot"), "System32", "schtasks.exe"))
		quoted := make([]string, len(args))
		for i, arg := range args {
			quoted[i] = syscall.EscapeArg(arg)
		}
		parameters, _ := windows.UTF16PtrFromString(strings.Join(quoted, " "))
		if err := windows.ShellExecute(0, verb, file, parameters, nil, windows.SW_HIDE); err != nil {
			return current, fmt.Errorf("Windows startup change was cancelled or denied: %w", err)
		}
	}

	deadline := time.Now().Add(60 * time.Second)
	for {
		found, actual, err := queryTask()
		if err == nil && actual == enabled && (found || !enabled) {
			return actual, nil
		}
		if time.Now().After(deadline) {
			return current, errors.New("Windows startup did not change; check Task Scheduler permissions")
		}
		time.Sleep(250 * time.Millisecond)
	}
}
