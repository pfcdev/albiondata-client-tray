//go:build windows

// Package consoleviewer provides an optional, separate log console. Closing
// this console exits only the viewer process, never the tray application.
package consoleviewer

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	procAllocConsole    = kernel32.NewProc("AllocConsole")
	procSetConsoleTitle = kernel32.NewProc("SetConsoleTitleW")
)

// Run allocates a console and follows the application's log file until the
// user closes the console or the parent closes its shutdown pipe.
func Run(path string, shutdown io.Reader) error {
	if result, _, err := procAllocConsole.Call(); result == 0 {
		return fmt.Errorf("AllocConsole: %w", err)
	}
	console, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer console.Close()
	if title, err := syscall.UTF16PtrFromString("Albion Data Client - Console"); err == nil {
		procSetConsoleTitle.Call(uintptr(unsafe.Pointer(title)))
	}

	// The parent keeps this pipe open while the viewer is wanted. Its close
	// also makes the viewer exit if the tray process stops unexpectedly.
	go func() {
		_, _ = io.Copy(io.Discard, shutdown)
		os.Exit(0)
	}()

	const historyBytes = 64 * 1024
	var file *os.File
	var offset int64
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	fmt.Fprintln(console, "Albion Data Client live log. Close this window to return to the tray.")
	for {
		if file == nil {
			file, err = os.Open(path)
			if err == nil {
				info, statErr := file.Stat()
				if statErr == nil && info.Size() > historyBytes {
					offset = info.Size() - historyBytes
				}
			} else {
				<-ticker.C
				continue
			}
		}
		info, err := file.Stat()
		if err != nil {
			file.Close()
			file = nil
			offset = 0
		} else {
			if info.Size() < offset {
				offset = 0
			}
			if _, err = file.Seek(offset, io.SeekStart); err == nil {
				_, err = io.Copy(console, file)
				if err == nil {
					offset, err = file.Seek(0, io.SeekCurrent)
				}
			}
			if err != nil {
				return err
			}
		}
		<-ticker.C
	}
}
