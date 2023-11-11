package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var oauthConfig2 = oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Scopes:       []string{"public_repo", "read:user", "user:email", "user:follow"},
	RedirectURL:  "https://szmul-med.onrender.com/github_user",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func LoggedHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(githubData), "", "\t"); err != nil {
		log.Println("JSON parse error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func (h *handlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="/login/">LOGIN</a>`)
}

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := oauthConfig2.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal("OAuth exchange failed:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	githubData := getGithubData(token.AccessToken)

	var githubUser model.GithubUser
	if err := json.Unmarshal([]byte(githubData), &githubUser); err != nil {
		log.Println("Error parsing GitHub data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the user already exists in the database
	existingUser := model.GithubUser{}
	if err := h.db.Where("login = ?", githubUser.Login).First(&existingUser).Error; err == nil {
		// Update existing record
		existingUser.Email = githubUser.Email
		// Update other fields as needed
		err := h.db.Save(&existingUser).Error
		if err != nil {
			log.Println("Failed to update github user in database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Create new record if it doesn't exist
		err := h.db.Create(&githubUser).Error
		if err != nil {
			log.Println("Failed to save user to database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	var publicRepo model.PublicRepo
	if err := json.Unmarshal([]byte(githubData), &publicRepo); err != nil {
		log.Println("Error parsing GitHub data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Check if the public repository already exists in the database
	existingRepo := model.PublicRepo{}
	if err := h.db.Where("repo_name = ?", publicRepo.Name).First(&existingRepo).Error; err == nil {
		// Update existing record
		existingRepo.Description = publicRepo.Description
		// Update other fields as needed
		err := h.db.Save(&existingRepo).Error
		if err != nil {
			log.Println("Failed to update public repository in database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else {
		// Create new record if it doesn't exist
		err := h.db.Create(&publicRepo).Error
		if err != nil {
			log.Println("Failed to save public repository to database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	hasAdminAccess := checkRepoAdminAccess(githubUser.AccessToken, "https://github.com/szmulinho/szmul-med")

	if hasAdminAccess {
		githubUser.Role = "admin"
	} else {
		githubUser.Role = "user"
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

	// Send user data to the user-api for registration
	resp, err := http.Post("https://szmul-med-users.onrender.com/register", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		log.Println("Failed to create user in user-api:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	LoggedHandler(w, r, githubData)
}

func checkRepoAdminAccess(accessToken, repoURL string) bool {
	u, err := url.Parse(repoURL)
	if err != nil {
		return false
	}

	pathComponents := strings.Split(u.Path, "/")
	if len(pathComponents) < 3 {
		return false
	}

	owner, repoName := pathComponents[1], pathComponents[2]

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})))

	_, _, err = client.Repositories.Get(context.Background(), owner, repoName)
	return err == nil
}

func getGithubData(accessToken string) string {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		log.Println("API Request creation failed:", err)
		return ""
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Request failed:", err)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Response read failed:", err)
		return ""
	}

	return string(respBody)
}
