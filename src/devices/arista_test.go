package devices

import (
	"testing"
)

func TestParseResult(t *testing.T) {
	d := &DeviceArista{}

	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Matching input",
			input:    "!this is a test \nend",
			expected: "!this is a test \nend",
		},
		{
			name:     "No endings",
			input:    "!this is a test",
			expected: "",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Complex example",
			input:    "!123\n123 end \n!456\nend",
			expected: "!123\n123 end \n!456\nend",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := d.ParseResult(tc.input)
			if err != nil && err.Error() != "Arista configuration match not found" {
				t.Errorf("Got unexpected error: %v", err)
			}
			if output != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, output)
			}
		})
	}
}
