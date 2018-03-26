package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"log"
	"regexp"
)

// devCiscoCsb pulls the device config for a Cisco Small Business device.
type devCiscoCsb struct {
	configuration.DeviceConfig
}

func makeCiscoCsb(d configuration.DeviceConfig) Collector {
	return &devCiscoCsb{d}
}

// BuildBatcher for CiscoCSB
// This is assuming prompt for User Name on Cisco CSB - this may not always be the case
func (p *devCiscoCsb) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{"User Name:", p.DeviceConfig.User + "\n"},
		{"Password:", p.DeviceConfig.Pass + "\n"},
		{".*#", "terminal datadump\n"},
		{".*#", "show running-config\n"},
		{".*#"},
	})
}

// ParseResult for CiscoCSB
func (p *devCiscoCsb) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	// This may break if there is a '#' in the config
	matcher := regexp.MustCompile(`config-file-header[\s\S]*?#`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}

// ModifySSHConfig since CiscoCSB needs special ciphers.
func (p *devCiscoCsb) ModifySSHConfig(config *utils.SSHConfig) {
	log.Println("Including Ciphers for Cisco CSB.")

	config.Ciphers = []string{"aes256-cbc", "aes128-cbc"}
}
