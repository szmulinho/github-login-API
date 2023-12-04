package endpoints

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func (h *handlers) getData(accessToken, apiUrl string) (string, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Printf("API Request creation failed: %v", err)
		return "", fmt.Errorf("API Request creation failed: %w", err)
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Error closing response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-OK status code received: %d", resp.StatusCode)
		return "", fmt.Errorf("Non-OK status code received: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Response read failed: %v", err)
		return "", fmt.Errorf("response read failed: %w", err)
	}

	return string(respBody), nil
}