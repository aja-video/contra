package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// DeviceEdgeSwitch logic container for device.
type DeviceEdgeSwitch struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceEdgeSwitch) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for EdgeSwitch
func (p *DeviceEdgeSwitch) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{".*>", "enable"},
		{"Password:", p.UnlockPass},
		{".*#", "terminal length 0"},
		{".*#", "show run"},
		{`\(.*\) #`},
	})
}

// ParseResult for EdgeSwitch
func (p *DeviceEdgeSwitch) ParseResult(result string) (string, error) {

	matcher := regexp.MustCompile(`hostname[\s\S]*exit`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
