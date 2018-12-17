package collectors

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// deviceArista pulls the device config for an Arista device
type deviceArista struct {
	configuration.DeviceConfig
}

func makeArista(d configuration.DeviceConfig) Collector {
	return &deviceArista{d}
}

// BuildBatcher for Arista
func (p *deviceArista) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{".*>", "terminal length 0"},
		{".*>", "enable"},
		{"Password:", p.UnlockPass},
		{".*#", "show run"},
		{"end"},
	})
}

// ParseResult for Arista
func (p *deviceArista) ParseResult(result string) (string, error) {

	matcher := regexp.MustCompile(`![\s\S]*end`)
	match := matcher.FindStringSubmatch(result)
	return match[0], nil
}
