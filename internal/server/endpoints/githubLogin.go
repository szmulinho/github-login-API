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

var (
	oauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		Scopes: []string{
			"public_repo"},
		RedirectURL: "https://szmul-med.onrender.com/github_user",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func LoggedHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func (h *handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/github/login/">LOGIN</a>`)
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := oauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)

	http.Redirect(w, r, redirectURL, 301)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal("OAuth exchange failed:", err)
	}

	githubData := getGithubData(token.AccessToken)

	var githubUser model.GithubUser
	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		log.Panic("Error parsing GitHub data:", err)
	}

	user := model.GithubUser{
		ID:           githubUser.ID,
		Login:        githubUser.Login,
		AvatarUrl:    githubUser.AvatarUrl,
		HtmlUrl:      githubUser.HtmlUrl,
		Email:        githubUser.Email,
		Role:         githubUser.Role,
		AccessToken:  token.AccessToken,
		Repositories: githubUser.Repositories,
	}

	if err := h.db.Create(&user).Error; err != nil {
		log.Panic("Failed to save user to database:", err)
	}

	hasAdminAccess := checkRepoAdminAccess(githubUser.AccessToken, "https://github.com/szmulinho/szmul-med")

	if hasAdminAccess {
		githubUser.Role = "admin"
	} else {
		githubUser.Role = "user"
	}

	newUser := model.GithubUser{
		Login: githubUser.Login,
		Email: githubUser.Email,
		Role:  githubUser.Role,
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://szmul-med-users.onrender.com/register", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Println("Failed to create user in user-api:", err)
		http.Error(w, "Failed to create user in user-api", http.StatusInternalServerError)
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

	owner := pathComponents[1]
	repoName := pathComponents[2]

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})))

	_, _, err = client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		return false
	}

	return true
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
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
