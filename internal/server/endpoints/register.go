package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

func (h *handlers) Register(w http.ResponseWriter, r *http.Request) {
	var newGithubUser model.GithubUser
	var publicRepos []model.PublicRepo
	szmulMedRepoName := "szmul-med"
	registerAPIBaseURL := "https://szmul-med-users.onrender.com/register"
	registerAPIDoctorsURL := "https://szmul-med-doctors.onrender.com/register"

	err := json.NewDecoder(r.Body).Decode(&newGithubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userJSON, err := json.Marshal(newGithubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		newGithubUser.Role = "doctor"
	} else {
		newGithubUser.Role = "user"
	}

	registerAPIURL := registerAPIBaseURL
	if hasSzmulMedRepo {
		registerAPIURL = registerAPIDoctorsURL
	}

	resp, err := http.Post(registerAPIURL, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		handleError(w, "Failed to create user in user-api", http.StatusInternalServerError, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Print(err)
		}
	}(resp.Body)

	result := h.db.Create(&newGithubUser)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}



	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}

func handleError(w http.ResponseWriter, errMsg string, status int, err error) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, status)
}
