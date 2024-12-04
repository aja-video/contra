package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// DeviceArista logic container for device.
type DeviceArista struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceArista) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Arista
func (p *DeviceArista) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{".*>", "terminal length 0"},
		{".*>", "enable"},
		{".*#", "show run"},
		{"end"},
	})
}

// ParseResult for Arista
func (p *DeviceArista) ParseResult(result string) (string, error) {

	matcher := regexp.MustCompile(`![\s\S]*end`)
	match := matcher.FindStringSubmatch(result)
	return match[0], nil
}
