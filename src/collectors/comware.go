package collectors

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"log"
	"regexp"
)

type DeviceComware struct {
	configuration.DeviceConfig
}

// SetConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceComware) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Comware
func (p *DeviceComware) BuildBatcher() ([]expect.Batcher, error) {
	if len(p.UnlockPass) > 0 {
		return utils.SimpleBatcher([][]string{
			{"<.*.>", "xtd-cli-mode"},
			{`(\[Y\/N\]\:$)`, "Y"},
			{"Password:", p.UnlockPass},
			{"<.*.>", "screen-length disable"},
			{"<.*.>", "display current-configuration"},
			{"return"},
		})
	}
	return utils.SimpleBatcher([][]string{
		{"<.*.>", "screen-length disable"},
		{"<.*.>", "display current-configuration"},
		{"return"},
	})
}

// ParseResult for Comware
func (p *DeviceComware) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`#[\s\S]*?return`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}

// ModifySSHConfig to add ciphers for locked down comware devices - Aruba 1950 for example
func (p *DeviceComware) ModifySSHConfig(config *utils.SSHConfig) {
	if len(p.UnlockPass) > 0 {
		log.Println("Including ciphers for comware with xtd-cli-mode")
		config.Ciphers = []string{"aes128-cbc", "aes256-cbc", "3des-cbc", "des-cbc"}
	}
}
