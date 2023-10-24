package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	oauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		RedirectURL:  "https://szmul-med.onrender.com/callback",
		Scopes:       []string{"user:email", "repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	// Set return type JSON
	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), http.StatusFound)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error exchanging code for token:", err)
		return
	}

	log.Println("Access Token:", token.AccessToken)

	githubUser := getUserInfoFromGitHub(token.AccessToken)
	githubData := getGithubData(token.AccessToken)
	log.Println("GitHub User Info:", githubUser)
	h.db.Create(&githubUser)

	http.Redirect(w, r, "/success", http.StatusFound)
	loggedinHandler(w, r, githubData)
}

func getUserInfoFromGitHub(accessToken string) model.GithubUser {
	client := oauthConfig.Client(context.Background(), &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Println("Error getting user info from GitHub:", err)
		return model.GithubUser{} // Obsłuż błąd i zwróć odpowiednią wartość
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	var githubUser model.GithubUser
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		fmt.Println("Error decoding user info:", err)
		return model.GithubUser{} // Obsłuż błąd i zwróć odpowiednią wartość
	}

	return githubUser
}

func getGithubData(accessToken string) string {
	// Get request to a set URL
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	// Set the Authorization header before sending the request
	// Authorization: token XXXXXXXXXXXXXXXXXXXXXXXXXXX
	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	// Make the request
	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	// Read the response as a byte slice
	respbody, _ := ioutil.ReadAll(resp.Body)

	// Convert byte slice to string and return
	return string(respbody)
}
