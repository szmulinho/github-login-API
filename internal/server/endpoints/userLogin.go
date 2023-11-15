package endpoints

//
//import (
//	"encoding/json"
//	"github.com/szmulinho/github-login/internal/model"
//	"net/http"
//)
//
//func (h *handlers) Login(w http.ResponseWriter, r *http.Request) {
//	var credentials struct {
//		Login string `json:"login"`
//	}
//
//	err := json.NewDecoder(r.Body).Decode(&credentials)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	var githubUser model.GithubUser
//	result := h.db.Where("login = ?", credentials.Login).First(&githubUser)
//	if result.Error != nil {
//		http.Error(w, "Invalid login", http.StatusUnauthorized)
//		return
//	}
//
//	var isGithubUser bool
//
//	if githubUser.Role == "user" || githubUser.Role == "doctor" {
//		isGithubUser = true
//	}
//
//	_, err = h.GenerateToken(w, r, githubUser.Login, isGithubUser)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	response := model.LoginResponse{
//		GithubUser: githubUser,
//	}
//
//	responseJSON, err := json.Marshal(response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	w.Write(responseJSON)
//}
