package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"net/http"
)

var (
	oauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		Scopes:       []string{"user", "repo"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://github.com/login/oauth/authorize",
			TokenURL: "https://github.com/login/oauth/access_token",
		},
	}
)

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline), http.StatusFound)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		// Handle error
		return
	}

	client := github.NewClient(oauthConfig.Client(ctx, token))

	// Get user info
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		// Handle error
		return
	}

	fmt.Println("User's Name:", user.GetName())

	repos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{})
	if err != nil {
		// Handle error
		return
	}

	desiredRepos := []string{"szmulinho/szmul-med", "szmulinho/drugstore", "szmulinho/prescription"}

	var isAdmin bool
	for _, repo := range repos {
		for _, desiredRepo := range desiredRepos {
			if repo.GetFullName() == desiredRepo {
				isAdmin = true
				break
			}
		}
		if isAdmin {
			break
		}
	}
	var role string
	if isAdmin {
		role = "admin"
	} else {
		role = "user"
	}

	githubUser := model.GithubUser{
		ID:    user.GetID(),
		Name:  user.GetName(),
		Email: user.GetEmail(),
		Role:  role,
	}

	jsonResponse, err := json.Marshal(githubUser)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
