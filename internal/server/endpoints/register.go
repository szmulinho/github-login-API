package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

func (h *handlers) Register(w http.ResponseWriter, r *http.Request) {
	var newGithubUser model.GithubUser

	err := json.NewDecoder(r.Body).Decode(&newGithubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.db.Create(&newGithubUser)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	userJSON, err := json.Marshal(newGithubUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}
