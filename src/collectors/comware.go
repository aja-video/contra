package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"fmt"
	"github.com/google/goexpect"
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
	switch {
	case len(p.UnlockPass) > 0:
		return utils.SimpleBatcher([][]string{
			{"<.*.>", "xtd-cli-mode"},
			{`(\[Y\/N\]\:$)`, "Y"},
			{"Password:", p.UnlockPass},
			{"<.*.>", "screen-length disable"},
			{"<.*.>", "display current-configuration"},
			{"return"},
		})
	default:
		return utils.SimpleBatcher([][]string{
			{"<.*.>", "screen-length disable"},
			{"<.*.>", "display current-configuration"},
			{"return"},
		})
	}
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
		fmt.Println("Including ciphers for comware with xtd cli")
		config.Ciphers = []string{"aes128-cbc", "aes256-cbc", "3des-cbc", "des-cbc"}
	}
}
