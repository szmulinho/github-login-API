package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
)

const githubOAuthURL = "https://github.com/login/oauth/access_token"

func (h *handlers) GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing GitHub code", http.StatusBadRequest)
		return
	}

	// Exchange GitHub code for access token
	token, err := h.exchangeGitHubCodeForToken(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange GitHub code for token", http.StatusInternalServerError)
		return
	}

	// Use the access token to get user information from GitHub
	githubUser, err := h.getGitHubUserInfo(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to get GitHub user information", http.StatusInternalServerError)
		return
	}

	// Here, you can save the GitHub user information to your database
	// For example:
	newUser := model.GithubUser{
		Username: githubUser.Username,
		Email:    githubUser.Email,
	}
	h.db.Create(&newUser)

	// Return the GitHub user information in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(githubUser)
}

func (h *handlers) exchangeGitHubCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	// Implement code to exchange GitHub code for access token using OAuth2
	// You can use the golang.org/x/oauth2 package for this purpose.
	// Example:
	config := oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		RedirectURL:  "your-redirect-url",
		Endpoint: oauth2.Endpoint{
			TokenURL: githubOAuthURL,
		},
	}
	token, err := config.Exchange(ctx, code)
	return token, err
	return nil, nil
}

func (h *handlers) getGitHubUserInfo(ctx context.Context, token *oauth2.Token) (*model.GithubUser, error) {
	// Implement code to get user information from GitHub using the access token
	// You can make a GET request to the GitHub API endpoint to get user information.
	// Example:
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var githubUser model.GithubUser
	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	return &githubUser, err
	return nil, nil
}

func (h *handlers) GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing GitHub code", http.StatusBadRequest)
		return
	}

	// Exchange GitHub code for access token
	token, err := h.exchangeGitHubCodeForToken(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange GitHub code for token", http.StatusInternalServerError)
		return
	}

	// Use the access token to get user information from GitHub
	githubUser, err := h.getGitHubUserInfo(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to get GitHub user information", http.StatusInternalServerError)
		return
	}

	// Here, you can save the GitHub user information to your database
	// For example:
	newUser := model.GithubUser{
		Username: githubUser.Username,
		Email:    githubUser.Email,
	}
	h.db.Create(&newUser)

	// Return the GitHub user information in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(githubUser)
}
