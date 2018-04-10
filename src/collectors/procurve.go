package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type deviceProcurve struct {
	configuration.DeviceConfig
}

func makeProcurve(d configuration.DeviceConfig) Collector {
	return &deviceProcurve{d}
}

// BuildBatcher for Procurve
func (p *deviceProcurve) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"continue", "a\n"},
		{".*#", "no page\n"},
		{".*#", "show running-config\n"},
		{".*[\\S]#"},
	})
}

// ParseResult for Procurve
func (p *deviceProcurve) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	// this regex assumes all procurve configs begin with 'hostname', and end with 'password manager'
	// Should probably find a better match...
	matcher := regexp.MustCompile(`hostname[\s\S]*?manager`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
