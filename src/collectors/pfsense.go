package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type devPfsense struct {
	configuration.DeviceConfig
}

func makePfsense(d configuration.DeviceConfig) Collector {
	return &devPfsense{d}
}

// BuildBatcher for pfSense
// "option:" should always match the initial connection string
func (p *devPfsense) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"option:", "8\n"},
		{".*root", "cat /conf/config.xml\n"},
		{"</pfsense>"},
	})
}

// ParseResult for pfSense
func (p *devPfsense) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
