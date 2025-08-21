package devices

import (
	"testing"
)

func TestParseResultNetgear(t *testing.T) {
	d := &DeviceNetgear{}

	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Matching input",
			input:    "!this\nconfigure is a test \n\n(M4500) #",
			expected: "configure is a test \n\n(M4500) #",
		},
		{
			name:     "No endings",
			input:    "!this\nconfigure is a test",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Complex example",
			input:    "!123\nconfigure\n123 'hello'\n!456\n\n(system-name) #",
			expected: "configure\n123 'hello'\n!456\n\n(system-name) #",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := d.ParseResult(tc.input)
			if err != nil && err.Error() != "netgear configuration match not found" {
				t.Errorf("Got unexpected error: %v", err)
			}
			if output != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, output)
			}
		})
	}
}
