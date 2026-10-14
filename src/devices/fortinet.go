package devices

import (
	"fmt"
	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	"github.com/google/goexpect"
	"regexp"
)

// DeviceFortinet logic container for device.
type DeviceFortinet struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceFortinet) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for Fortinet
func (p *DeviceFortinet) BuildBatcher() ([]expect.Batcher, error) {
	return utils.SimpleBatcher([][]string{
		{".* #", "config system console"},
		{".*\\(console\\) #", "set output standard"},
		{".*\\(console\\) #", "end"},
		{".*#", "show"},
		{"\n.* # $"},
	})
}

// ParseResult for Fortinet
func (p *DeviceFortinet) ParseResult(result string) (string, error) {
	matcher := regexp.MustCompile(`config system global[\s\S]*end`)
	match := matcher.FindStringSubmatch(result)
	if len(match) == 0 {
		return "", fmt.Errorf("fortinet configuration match not found")
	}
	// redact encrypted passwords as they change every run
	// base64 padding varies from zero to two '=', so match the whole whitespace-free value
	// anchor on FortiOS's `set <attribute> ENC <value>` form so "ENC " inside a quoted
	// object name or comment is left alone
	encryptedPassword := regexp.MustCompile(`(?m)^([ \t]*set \S+ ENC )\S+`)
	redactedResult := encryptedPassword.ReplaceAllString(match[0], "${1}~~~Contra Redacted~~~")

	return redactedResult, nil
}
