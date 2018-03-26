package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type devComware struct {
	configuration.DeviceConfig
}

func makeComware(d configuration.DeviceConfig) Collector {
	return &devComware{d}
}

// BuildBatcher for Comware
func (p *devComware) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"<.*.>", "screen-length disable\n"},
		{"<.*.>", "display current-configuration\n"},
		{"return"},
	})
}

// ParseResult for Comware
func (p *devComware) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`#[\s\S]*?return`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
