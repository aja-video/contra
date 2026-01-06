package devices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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

func doHTTPRequest(p *DeviceHTTPJSON) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	uri := "https://" + p.Host + ":" + strconv.Itoa(p.Port) + p.Endpoint

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Token "+p.Pass)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return body, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return body, nil
}

// ParseResult for http
func (p *DeviceHTTPJSON) ParseResult(result string) (string, error) {
	body, err := doHTTPRequest(p)
	if err != nil {
		return string(body), err
	}

	return formatJSON(body)
}
