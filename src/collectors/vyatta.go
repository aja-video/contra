package collectors

import (
	"contra/src/configuration"
	"contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// deviceVyatta pulls the device config for a Vyatta based device.
type deviceVyatta struct {
	configuration.DeviceConfig
}

func makeVyatta(d configuration.DeviceConfig) Collector {
	return &deviceVyatta{d}
}

// BuildBatcher for Vyatta
func (p *deviceVyatta) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{`.*\$`, "terminal length 0" + "\n"},
		{`.*\$`, "show configuration" + "\n"},
		{`.*\$`},
	})
}

// ParseResult for Vyatta
func (p *deviceVyatta) ParseResult(result string) (string, error) {
	matcher := regexp.MustCompile(`(.*\{[\s\S]*\})\n[\S\s]*\$`)
	match := matcher.FindStringSubmatch(result)
	return match[1], nil
}