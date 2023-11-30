package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

func (h *handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser model.User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser.Role = "user"

	result := h.db.Create(&newUser)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userJSON)
}
