package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"log"
	"regexp"
)

type deviceComware struct {
	configuration.DeviceConfig
}

func makeComware(d configuration.DeviceConfig) Collector {
	return &deviceComware{d}
}

// BuildBatcher for Comware
func (p *deviceComware) BuildBatcher() ([]expect.Batcher, error) {
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
func (p *deviceComware) ParseResult(result string) (string, error) {
	// Strip shell commands, grab only the xml file
	matcher := regexp.MustCompile(`#[\s\S]*?return`)
	match := matcher.FindStringSubmatch(result)

	return match[0], nil
}

// ModifySSHConfig to add ciphers for locked down comware devices - Aruba 1950 for example
func (p *deviceComware) ModifySSHConfig(config *utils.SSHConfig) {
	if len(p.UnlockPass) > 0 {
		log.Println("Including ciphers for comware with xtd-cli-mode")
		config.Ciphers = []string{"aes128-cbc", "aes256-cbc", "3des-cbc", "des-cbc"}
	}
}
