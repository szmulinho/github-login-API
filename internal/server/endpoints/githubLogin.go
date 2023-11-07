package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var oauthConfig = oauth2.Config{
	ClientID:     "33f5f8298ded51f76f30",
	ClientSecret: "1b7ab1c0faeac3b5b3619bfe610efc9514713f85",
	Scopes:       []string{"public_repo", "read:user", "user:email", "user:follow", "read:project", "read:packages"},
	RedirectURL:  "https://szmul-med.onrender.com/github_user",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func LoggedHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(githubData), "", "\t"); err != nil {
		log.Println("JSON parse error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func (h *handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/github/login/">LOGIN</a>`)
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := oauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal("OAuth exchange failed:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	githubData := getGithubData(token.AccessToken)
	var githubUser model.GithubUser
	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		log.Println("Error parsing GitHub data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hasAdminAccess := checkRepoAdminAccess(githubUser.AccessToken, "https://github.com/szmulinho/szmul-med")

	if hasAdminAccess {
		githubUser.Role = "admin"
	} else {
		githubUser.Role = "user"
	}

	// Save user to the database
	err = h.db.Create(&githubUser).Error
	if err != nil {
		log.Println("Failed to save user to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a simplified user object for response
	newUser := model.GithubUser{
		Login: githubUser.Login,
		Email: githubUser.Email,
		Role:  githubUser.Role,
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		log.Println("JSON marshaling error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send user data to the user-api for registration
	resp, err := http.Post("https://szmul-med-users.onrender.com/register", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Println("Failed to create user in user-api:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	LoggedHandler(w, r, githubData)
}

func checkRepoAdminAccess(accessToken, repoURL string) bool {
	u, err := url.Parse(repoURL)
	if err != nil {
		return false
	}

	pathComponents := strings.Split(u.Path, "/")
	if len(pathComponents) < 3 {
		return false
	}

	owner, repoName := pathComponents[1], pathComponents[2]

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})))

	_, _, err = client.Repositories.Get(context.Background(), owner, repoName)
	return err == nil
}

func getGithubData(accessToken string) string {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Println("API Request creation failed:", err)
		return ""
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request failed:", err)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Response read failed:", err)
		return ""
	}

	return string(respBody)
}
