package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// DeviceVyatta logic container for device.
type DeviceVyatta struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceVyatta) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Vyatta
func (p *DeviceVyatta) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{`.*\$`, "terminal length 0"},
		{`.*\$`, "show configuration"},
		{`.*\$`},
	})
}

// ParseResult for Vyatta
func (p *DeviceVyatta) ParseResult(result string) (string, error) {
	matcher := regexp.MustCompile(`(.*\{[\s\S]*\})\n[\S\s]*\$`)
	match := matcher.FindStringSubmatch(result)
	return match[1], nil
}
