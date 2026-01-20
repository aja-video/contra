package devices

import (
	"fmt"
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	expect "github.com/google/goexpect"
	"regexp"
)

// DeviceNetgear logic container for device.
type DeviceNetgear struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceNetgear) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Netgear
func (p *DeviceNetgear) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{".*#", "no pager"},
		{".*#", "show run"},
		{".*#"},
	})
}

// ParseResult for Netgear
func (p *DeviceNetgear) ParseResult(result string) (string, error) {

	matcher := regexp.MustCompile(`(?ms)configure.*?#$`)
	match := matcher.FindStringSubmatch(result)
	if len(match) == 0 {
		return "", fmt.Errorf("netgear configuration match not found")
	}
	return match[0], nil
}
