package endpoints

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func (h *handlers)getUserData(accessToken string, endpoint string) string {

	req, reqerr := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.github.com/user/%s", endpoint),
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}
