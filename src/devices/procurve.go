package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// DeviceProcurve logic container for device.
type DeviceProcurve struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceProcurve) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Procurve
func (p *DeviceProcurve) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"continue", "a"},
		{".*#", "no page"},
		{".*#", "show running-config"},
		{".*[\\S]#"},
	})
}

// ParseResult for Procurve
func (p *DeviceProcurve) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	// this regex assumes all procurve configs begin with 'hostname', and end with 'password manager'
	// Should probably find a better match...
	matcher := regexp.MustCompile(`hostname[\s\S]*?manager`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
