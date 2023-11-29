package endpoints

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

const (
	szmulMedRepoName      = "szmul-med"
	registerAPIBaseURL    = "https://szmul-med-users.onrender.com/register"
	registerAPIDoctorsURL = "https://szmul-med-doctors.onrender.com/register"
)

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig2.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusBadRequest)
		log.Println("Error exchanging code for token:", err)
		return
	}

	var githubUser model.GithubUser

	userURL := "https://api.github.com/user"
	githubData, err := h.getData(token.AccessToken, userURL)
	if err != nil {
		handleError(w, "Error fetching user data from GitHub", http.StatusInternalServerError, err)
		return
	}

	response := model.LoginResponse{
		GithubUser: githubUser,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)


	reposURL := "https://api.github.com/user/repos"
	reposResp, err := h.getData(token.AccessToken, reposURL)
	if err != nil {
		handleError(w, "Error fetching user repositories", http.StatusInternalServerError, err)
		return
	}

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

	registerAPIURL := registerAPIBaseURL
	if hasSzmulMedRepo {
		registerAPIURL = registerAPIDoctorsURL
	}

	newUser := model.GithubUser{
		Login: githubUser.Login,
		Email: githubUser.Email,
		Role:  githubUser.Role,
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		handleError(w, "JSON marshaling error", http.StatusInternalServerError, err)
		return
	}

	resp, err := http.Post(registerAPIURL, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		handleError(w, "Failed to create user in user-api", http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	h.Logged(w, r, githubData)
}

func handleError(w http.ResponseWriter, errMsg string, status int, err error) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, status)
}
