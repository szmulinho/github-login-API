package endpoints

import (
	"encoding/json"
	"github.com/szmulinho/github-login/internal/model"
	"net/http"
)

func (h *handlers) GetUserData(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Query().Get("login")
	if login == "" {
		http.Error(w, "Missing user login parameter", http.StatusBadRequest)
		return
	}

	var userData model.GithubUser

	if err := h.db.Where("login = ?", login).First(&userData).Error; err != nil {
		http.Error(w, "Error fetching user data from the database", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Error marshaling user data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
