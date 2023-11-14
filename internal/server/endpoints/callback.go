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

	reposURL := "https://api.github.com/user/repos"
	reposResp, err := h.getData(token.AccessToken, reposURL)
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
	if err := json.Unmarshal([]byte(reposResp), &publicRepos); err != nil {
		log.Println("Error parsing GitHub repositories data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hasSzmulMedRepo := false
	for _, repo := range publicRepos {
		if repo.Name == "szmul-med" {
			hasSzmulMedRepo = true
			break
		}
	}

	if hasSzmulMedRepo {
		githubUser.Role = "doctor"
	} else {
		githubUser.Role = "user"
	}

	existingUser := model.GithubUser{}
	if err := h.db.Where("login = ?", githubUser.Login).First(&existingUser).Error; err == nil {
		existingUser.Email = githubUser.Email
		err := h.db.Save(&existingUser).Error
		if err != nil {
			log.Println("Failed to update github user in database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		err := h.db.Create(&githubUser).Error
		if err != nil {
			log.Println("Failed to save user to database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	token2, err := h.GenerateToken(w, r, githubUser, true)
	if err != nil {
		// Handle error, perhaps return an internal server error
		return
	}

	var publicRepo model.PublicRepo
	if err := json.Unmarshal([]byte(githubData), &publicRepo); err != nil {
		log.Println("Error parsing GitHub data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	existingRepo := model.PublicRepo{}
	if err := h.db.Where("name = ?", publicRepo.Name).First(&existingRepo).Error; err == nil {
		existingRepo.Description = publicRepo.Description
		err := h.db.Save(&existingRepo).Error
		if err != nil {
			log.Println("Failed to update public repository in database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		err := h.db.Create(&publicRepo).Error
		if err != nil {
			log.Println("Failed to save public repository to database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	registerAPIURL := "https://szmul-med-users.onrender.com/register"
	if hasSzmulMedRepo {
		registerAPIURL = "https://szmul-med-doctors.onrender.com/register"
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

	resp, err := http.Post(registerAPIURL, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Println("Failed to create user in user-api:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	h.getUserFromToken(token2)

	Logged(w, r, githubData)
}
