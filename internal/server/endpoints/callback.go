package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

const (
	szmulMedRepoName      = "szmul-med"
)


func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := Config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusBadRequest)
		log.Println("Error exchanging code for token:", err)
		return
	}

	userURL := "https://api.github.com/user"
	githubData, err := h.getData(token.AccessToken, userURL)
	if err != nil {
		handleError(w, "Error fetching user data from GitHub", http.StatusInternalServerError, err)
		return
	}

	reposURL := "https://api.github.com/user/repos"
	reposResp, err := h.getData(token.AccessToken, reposURL)
	if err != nil {
		handleError(w, "Error fetching user repositories", http.StatusInternalServerError, err)
		return
	}

	var githubUser model.GithubUser
	var publicRepos []model.PublicRepo
	var publicRepo model.PublicRepo

	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		handleError(w, "Error parsing GitHub data", http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal([]byte(reposResp), &publicRepos); err != nil {
		handleError(w, "Error parsing GitHub repositories data", http.StatusInternalServerError, err)
		return
	}

	var hasSzmulMedRepo bool
	for _, repo := range publicRepos {
		if repo.Name == szmulMedRepoName {
			hasSzmulMedRepo = true
			break
		}
	}

	if hasSzmulMedRepo {
		githubUser.Role = "doctor"
	} else {
		githubUser.Role = "user"
	}


	if err := h.updateOrCreateGitHubUser(h.db, githubUser); err != nil {
		handleError(w, "Failed to update/create GitHub user in the database", http.StatusInternalServerError, err)
		return
	}

	if err := h.updateOrCreatePublicRepo(h.db, publicRepo); err != nil {
		handleError(w, "Failed to update/create public repository in the database", http.StatusInternalServerError, err)
		return
	}

	updatedGithubData, err := json.Marshal(githubUser)
	if err != nil {
		handleError(w, "Error serializing updated GitHub data", http.StatusInternalServerError, err)
		return
	}

	githubData = string(updatedGithubData)

	h.Logged(w, r, githubData)
}