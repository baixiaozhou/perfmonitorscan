package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendData(endpoint string, data interface{}) error {
	// json data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// Create Http request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
