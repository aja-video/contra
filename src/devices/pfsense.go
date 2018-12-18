package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

type DevicePfsense struct {
	configuration.DeviceConfig
}

// SetConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DevicePfsense) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for pfSense
func (p *DevicePfsense) BuildBatcher() ([]expect.Batcher, error) {
	// The "OK" result must be the first entry for variable.
	// The more the better, since this is constantly checking every case for a match.
	// - So simply having .*root will match multiple times throughout the dump.
	return utils.VariableBatcher([][]string{
		{`</pfsense>`}, // Found the "OK" result!
		{`Enter an option: `, "8"},
		{`/root.*:`, "cat /conf/config.xml"},
		{`\$ `, "cat /conf/config.xml"},
	})
}

// ParseResult for pfSense
func (p *DevicePfsense) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`<\?xml version[\s\S]*?<\/pfsense>`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}
