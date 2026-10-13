package devices

import (
	"testing"
)

func TestParseResultFortinet(t *testing.T) {
	d := &DeviceFortinet{}

	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No padding",
			input:    "config system global\r\n    set password ENC vlSNIIw2J3ufy65wt8/wBylMq\r\nend",
			expected: "config system global\r\n    set password ENC ~~~Contra Redacted~~~\r\nend",
		},
		{
			name:     "Single padding",
			input:    "config system global\n    set password ENC c7fMh14XfSNZEKLyy=\nend",
			expected: "config system global\n    set password ENC ~~~Contra Redacted~~~\nend",
		},
		{
			name:     "Double padding",
			input:    "config system global\n    set password ENC F4O6YMjrMmDSeMuCfEtnOXLhdXK==\nend",
			expected: "config system global\n    set password ENC ~~~Contra Redacted~~~\nend",
		},
		{
			name: "Multiple values and untouched lines",
			input: "config system global\n    set api-key ENC SGVsbG8=\n" +
				"    set comments \"This is the default\"\n    set password ENC V29ybGQ\nend",
			expected: "config system global\n    set api-key ENC ~~~Contra Redacted~~~\n" +
				"    set comments \"This is the default\"\n    set password ENC ~~~Contra Redacted~~~\nend",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := d.ParseResult(tc.input)
			if err != nil {
				t.Errorf("Got unexpected error: %v", err)
			}
			if output != tc.expected {
				t.Errorf("Expected %q, but got %q", tc.expected, output)
			}
		})
	}
}
