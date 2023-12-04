package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

var githubUser model.GhUser
var publicRepos []model.PublicRepo

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusBadRequest)
		log.Println("Error exchanging code for token:", err)
		return
	}

	userURL := "https://api.github.com/user"
	userData, err := getData(token.AccessToken, userURL)
	if err != nil {
		http.Error(w, "Error fetching user data from GitHub", http.StatusInternalServerError)
		return
	}

	reposURL := "https://api.github.com/user/repos"
	reposData, err := getData(token.AccessToken, reposURL)
	if err != nil {
		http.Error(w, "Error fetching user repositories", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal([]byte(userData), &githubUser); err != nil {
		http.Error(w, "Error parsing GitHub data", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal([]byte(reposData), &publicRepos); err != nil {
		http.Error(w, "Error parsing GitHub repositories data", http.StatusInternalServerError)
		return
	}

	var hasSzmulMedRepo bool
	for _, repo := range publicRepos {
		if repo.Name == "szmul-med" {
			hasSzmulMedRepo = true
			break
		}
	}

	if hasSzmulMedRepo {
		githubUser.Role = "doctor"
	} else {
		githubUser.Role = "user"
	}


	if err := updateOrCreateGitHubUser(h.db, githubUser); err != nil {
		http.Error(w, "Failed to update/create GitHub user in the database", http.StatusInternalServerError)
		return
	}

	updatedGithubData, err := json.Marshal(githubUser)
	if err != nil {
		http.Error(w, "Error serializing updated GitHub data", http.StatusInternalServerError)
		return
	}

	userData = string(updatedGithubData)

	h.Logged(w, r, userData)
}