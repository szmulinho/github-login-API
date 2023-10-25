package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func LoggedHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		_, err := fmt.Fprintf(w, "UNAUTHORIZED!")
		if err != nil {
			return
		}
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parser := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parser != nil {
		log.Panic("JSON parse error")
	}

	_, err := fmt.Fprintf(w, string(prettyJSON.Bytes()))
	if err != nil {
		return
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, `<a href="/login/github/">LOGIN</a>`)
	if err != nil {
		return
	}
}

func (h *handlers) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Get the environment variable
	githubClientID := "065d047663d40d183c04"

	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		githubClientID,
		"https://szmul-med-github-login.onrender.com/login/github/callback",
	)

	http.Redirect(w, r, redirectURL, 301)
}

func (h *handlers) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	LoggedHandler(w, r, githubData)
}

func getGithubAccessToken(code string) string {

	clientID := "065d047663d40d183c04"
	clientSecret := "7b7c2239b98e0b66d53e6b2adbfd8722561512f4"

	// Set us the request body as JSON
	requestBodyMap := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// POST request to set URL
	req, reqerr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Get the response
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	err := json.Unmarshal(respbody, &ghresp)
	if err != nil {
		return ""
	}

	return ghresp.AccessToken
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	body, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(body)
}
