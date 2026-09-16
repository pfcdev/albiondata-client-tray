package autostart

import (
	"testing"
	"unicode/utf16"
)

func TestEnabledFromTaskXML(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want bool
	}{
		{"default enabled", `<Task><Settings><MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy></Settings></Task>`, true},
		{"disabled", `<Task><Settings><Enabled>false</Enabled></Settings></Task>`, false},
		{"explicit enabled", `<Task><Settings><Enabled>true</Enabled></Settings></Task>`, true},
		{"trigger disabled only", `<Task><Triggers><LogonTrigger><Enabled>false</Enabled></LogonTrigger></Triggers><Settings></Settings></Task>`, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := enabledFromTaskXML([]byte(tc.xml))
			if err != nil || got != tc.want {
				t.Fatalf("enabledFromTaskXML() = %v, %v; want %v", got, err, tc.want)
			}
		})
	}
}

func TestEnabledFromTaskXMLRejectsMissingSettings(t *testing.T) {
	if _, err := enabledFromTaskXML([]byte(`<Task/>`)); err == nil {
		t.Fatal("expected an error for missing task settings")
	}
}

func TestEnabledFromTaskXMLUTF16(t *testing.T) {
	text := `<Task><Settings><Enabled>false</Enabled></Settings></Task>`
	data := []byte{0xff, 0xfe}
	for _, unit := range utf16.Encode([]rune(text)) {
		data = append(data, byte(unit), byte(unit>>8))
	}
	enabled, err := enabledFromTaskXML(data)
	if err != nil || enabled {
		t.Fatalf("enabledFromTaskXML(UTF-16) = %v, %v; want false", enabled, err)
	}
}
