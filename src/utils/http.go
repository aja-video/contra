package utils

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func HTTPRequest(host string, port string, endpoint string, pass string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	uri := "https://" + host + ":" + port + endpoint

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	if pass != "" {
		req.Header.Add("Authorization", "Token "+pass)
	}

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
