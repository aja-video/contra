package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"log"
	"regexp"
)

// DeviceMikrotik logic container for device.
type DeviceMikrotik struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceMikrotik) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Mikrotik
func (p *DeviceMikrotik) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{`\[.*@.*] >`, "/export\r"},
		{`\/export`},
		{`\[.*@.*] >`},
	})
}

// ParseResult for Mikrotik
func (p *DeviceMikrotik) ParseResult(result string) (string, error) {
	// start match at software id - date stamp changes every time /export is run
	matcher := regexp.MustCompile(`(\# software id[\s\S]*?)?\[.*@.*] >`)
	match := matcher.FindStringSubmatch(result)

	return match[1], nil
}

// ModifyUsername to fix mikrotik terminal
func (p *DeviceMikrotik) ModifyUsername(config *utils.SSHConfig) {
	log.Printf("Updating %s terminal preferences", p.Name)

	config.User = config.User + "+ct"
}
