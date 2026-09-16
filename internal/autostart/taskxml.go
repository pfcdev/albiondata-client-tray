package autostart

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"unicode/utf16"
)

var enabledTag = regexp.MustCompile(`(?i)<Enabled>\s*(true|false)\s*</Enabled>`)

// enabledFromTaskXML reads the task-level setting, not an individual
// trigger's Enabled element. Task Scheduler omits it when the task is on.
func enabledFromTaskXML(data []byte) (bool, error) {
	if len(data) >= 2 && data[0] == 0xff && data[1] == 0xfe {
		if (len(data)-2)%2 != 0 {
			return false, errors.New("invalid UTF-16 task XML")
		}
		units := make([]uint16, (len(data)-2)/2)
		for i := range units {
			units[i] = uint16(data[2+i*2]) | uint16(data[3+i*2])<<8
		}
		data = []byte(string(utf16.Decode(units)))
	}
	start := bytes.Index(data, []byte("<Settings>"))
	if start < 0 {
		return false, errors.New("task XML has no Settings element")
	}
	settings := data[start+len("<Settings>"):]
	end := bytes.Index(settings, []byte("</Settings>"))
	if end < 0 {
		return false, errors.New("task XML has no closing Settings element")
	}
	match := enabledTag.FindSubmatch(settings[:end])
	if match == nil {
		return true, nil
	}
	return strings.EqualFold(string(match[1]), "true"), nil
}
