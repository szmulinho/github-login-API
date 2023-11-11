package endpoints

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *handlers) getData(accessToken, apiUrl string) (string, error) {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Println("API Request creation failed:", err)
		return "", err
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-OK status code received:", resp.StatusCode)
		return "", fmt.Errorf("Non-OK status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Response read failed:", err)
		return "", err
	}

	return string(respBody), nil
}
