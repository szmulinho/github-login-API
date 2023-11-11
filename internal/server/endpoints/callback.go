package endpoints

import (
	"bytes"
	"encoding/json"
	"github.com/szmulinho/github-login/internal/model"
	"log"
	"net/http"
)

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig2.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal("OAuth exchange failed:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userURL := "https://api.github.com/user"
	githubData, err := h.getData(token.AccessToken, userURL)
	if err != nil {
		log.Println("Error fetching user repositories:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var githubUser model.GithubUser
	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		log.Println("Error parsing GitHub data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var publicRepos []model.PublicRepo
	if err := json.Unmarshal([]byte(githubData), &publicRepos); err != nil {
		log.Println("Error parsing GitHub repositories data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, repo := range publicRepos {
		existingRepo := model.PublicRepo{}
		if err := h.db.Where("name = ?", repo.Name).First(&existingRepo).Error; err == nil {
			existingRepo.Description = repo.Description
			err := h.db.Save(&existingRepo).Error
			if err != nil {
				log.Println("Failed to update public repository in database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			err := h.db.Create(&repo).Error
			if err != nil {
				log.Println("Failed to save public repository to database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	newUser := model.GithubUser{
		Login: githubUser.Login,
		Email: githubUser.Email,
		Role:  githubUser.Role,
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		log.Println("JSON marshaling error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://szmul-med-users.onrender.com/register", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Println("Failed to create user in user-api:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	Logged(w, r, githubData)
}
