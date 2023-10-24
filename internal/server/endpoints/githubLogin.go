package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var (
	oauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		RedirectURL:  "https://szmul-med.onrender.com/github_user",
		Scopes:       []string{"user:email", "repo"},
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
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error exchanging code for token:", err)
		return
	}

	log.Println("Access Token:", token.AccessToken)

	githubUser := h.GetUserInfoFromGitHub(token.AccessToken)
	log.Println("GitHub User Info:", githubUser)

	if err := h.db.Create(&githubUser).Error; err != nil {
		log.Println("Error saving user to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}

func (h *handlers) GetUserInfoFromGitHub(accessToken string) model.GithubUser {
	client := oauthConfig.Client(context.Background(), &oauth2.Token{AccessToken: accessToken})
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		fmt.Println("Error getting user info from GitHub:", err)
		return model.GithubUser{} // Obsłuż błąd i zwróć odpowiednią wartość
	}
	defer resp.Body.Close()

	// Dodaj logi dla odpowiedzi od GitHub API
	fmt.Println("GitHub API Response:", resp.Status)

	var githubUser model.GithubUser
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		fmt.Println("Error decoding user info:", err)
		return model.GithubUser{} // Obsłuż błąd i zwróć odpowiednią wartość
	}

	return githubUser
}
