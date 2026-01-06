package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aja-video/contra/src/configuration"
	"github.com/aja-video/contra/src/utils"
	expect "github.com/google/goexpect"
)

// DeviceHTTPJSON logic container for device.
type DeviceHTTPJSON struct {
	configuration.DeviceConfig
}

// SetDeviceConfig since it is unclear how to assign DeviceConfig via reflect.New
func (p *DeviceHTTPJSON) SetDeviceConfig(deviceConfig configuration.DeviceConfig) {
	p.DeviceConfig = deviceConfig
}

// BuildBatcher for http-json
func (p *DeviceHTTPJSON) BuildBatcher() ([]expect.Batcher, error) {
	// no expect neccesary so return a stub
	return utils.SimpleBatcher([][]string{})
}

func formatJSON(rawJSON []byte) (string, error) {
	if len(rawJSON) == 0 {
		return "", fmt.Errorf("empty response body")
	}

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, rawJSON, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format JSON: %w", err)
	}

	return prettyJSON.String(), nil
}

// ParseResult for http
func (p *DeviceHTTPJSON) ParseResult(result string) (string, error) {
	body, err := utils.HTTPRequest(p.Host, strconv.Itoa(p.Port), p.Endpoint, p.Pass)
	if err != nil {
		return string(body), err
	}

	return formatJSON(body)
}
