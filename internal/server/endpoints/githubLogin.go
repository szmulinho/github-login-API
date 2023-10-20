package endpoints

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"net/http"
)

var (
	githubOauthConfig = oauth2.Config{
		ClientID:     "065d047663d40d183c04",
		ClientSecret: "7b7c2239b98e0b66d53e6b2adbfd8722561512f4",
		RedirectURL:  "http://localhost:5173/profile",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}
)

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	handler := gologin.NewHandler("https://github.com/", "client_id", "client_secret")

	handler.HandleLogin(w, r)
}

func (h *endpoints.handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	user, err := handler.UserFromContext(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hasAccess := hasAccessToRepositories(user)

	role := "user"
	if hasAccess {
		role = "admin"
	}

	insertUser(user, role)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *handlers) HasAccessToRepositories(user *gologin.User) bool {
	repositories, err := getRepositories(user)
	if err != nil {
		fmt.Println(err)
		return false
	}

	for _, repository := range repositories {
		if repository.Permissions.Admin {
			return true
		}
	}

	return false
}

func (h *handlers) GetRepositories(user *gologin.User) ([]*gologin.Repository, error) {
	request := gologin.RepositoryListRequest{
		User: user.Username,
	}

	// Pobierz listę repozytoriów
	repositories, err := handler.ListRepositories(request)
	if err != nil {
		return nil, err
	}

	return repositories, nil
}

func insertUser(user *gologin.User, role string) error {
	// Utwórz połączenie z bazą danych
	db, err := dbConnect()
	if err != nil {
		return err
	}

	// Wygeneruj zapytanie SQL
	stmt, err := db.Prepare("INSERT INTO users (username, role) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	// Wykonaj zapytanie SQL
	_, err = stmt.Exec(user.Username, role)
	if err != nil {
		return err
	}

	return nil
}
