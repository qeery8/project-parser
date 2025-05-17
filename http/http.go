package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetAPIResponse(url string, result interface{}, headers map[string]string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %s", err)
	}

	if headers == nil {
		headers = make(map[string]string)
	}
	if _, ok := headers["User-Agent"]; !ok {
		headers["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0"
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = "application/json"
	}
	if _, ok := headers["Accept-Language"]; !ok {
		headers["Accept-Language"] = "en-US,en;q=0.9"
	}
	if _, ok := headers["Origin"]; !ok {
		headers["Origin"] = "https://es.wallapop.com"
	}
	if _, ok := headers["Referer"]; !ok {
		headers["Referer"] = "https://es.wallapop.com/"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make the HTTP request: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read the response body: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(body))
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %s", err)
	}

	return nil
}
