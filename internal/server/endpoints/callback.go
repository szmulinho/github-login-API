package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
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
	_, err = w.Write(responseJSON)
	if err != nil {
		return
	}


	reposURL := "https://api.github.com/user/repos"
	reposResp, err := h.getData(token.AccessToken, reposURL)
	if err != nil {
		handleError(w, "Error fetching user repositories", http.StatusInternalServerError, err)
		return
	}

	var publicRepos []model.PublicRepo

	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		handleError(w, "Error parsing GitHub data", http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal([]byte(reposResp), &publicRepos); err != nil {
		handleError(w, "Error parsing GitHub repositories data", http.StatusInternalServerError, err)
		return
	}

	h.Logged(w, r, githubData)
}

func handleError(w http.ResponseWriter, errMsg string, status int, err error) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, status)
}
