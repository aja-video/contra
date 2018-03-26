package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type devProcurve struct {
	configuration.DeviceConfig
}

func makeProcurve(d configuration.DeviceConfig) Collector {
	return &devProcurve{d}
}

// BuildBatcher for Procurve
func (p *devProcurve) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"continue", "a\n"},
		{".*#", "no page\n"},
		{".*#", "show running-config\n"},
		{".*#"},
	})
}

// ParseResult for Procurve
func (p *devProcurve) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	// this regex assumes all procurve configs begin with 'hostname', and end with 'password manager'
	// Should probably find a better match...
	matcher := regexp.MustCompile(`hostname[\s\S]*?manager`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
