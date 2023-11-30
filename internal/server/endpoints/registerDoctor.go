package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/szmulinho/github-login/internal/model"
)

func (h *handlers) RegisterDoctor(w http.ResponseWriter, r *http.Request) {
	var newDoctor model.Doctor

	err := json.NewDecoder(r.Body).Decode(&newDoctor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newDoctor.Role = "doctor"

	result := h.db.Create(&newDoctor)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	doctorJSON, err := json.Marshal(newDoctor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(doctorJSON)
}
