package devices

import (
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
	// redact encrypted passwords as they change every run
	encryptedPassword := regexp.MustCompile("ENC .*==")
	redactedResult := encryptedPassword.ReplaceAllString(match[0], "ENC ~~~Contra Redacted~~~")

	return redactedResult, nil
}
