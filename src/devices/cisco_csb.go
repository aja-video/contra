package devices

import (
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"log"
	"regexp"
)

type DeviceCiscoCsb struct {
	configuration.DeviceConfig
}

// SetConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceCiscoCsb) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for CiscoCSB
// This is assuming prompt for User Name on Cisco CSB - this may not always be the case
func (p *DeviceCiscoCsb) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"User Name:", p.User},
		{"Password:", p.Pass},
		{".*#", "terminal datadump"},
		{".*#", "show running-config"},
		{".*#"},
	})
}

// ParseResult for CiscoCSB
func (p *DeviceCiscoCsb) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	// This may break if there is a '#' in the config
	matcher := regexp.MustCompile(`show.*[\s\S]\n(.*[\s\S]*)\n[\S.]*#`)
	match := matcher.FindStringSubmatch(result)

	return match[1], nil
}

// ModifySSHConfig since CiscoCSB needs special ciphers.
func (p *DeviceCiscoCsb) ModifySSHConfig(config *utils.SSHConfig) {
	log.Println("Including Ciphers for Cisco CSB.")

	config.Ciphers = []string{"aes256-cbc", "aes128-cbc"}
}
