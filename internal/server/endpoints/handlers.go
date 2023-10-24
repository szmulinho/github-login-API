package endpoints

import (
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type Handlers interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	GetUserInfoFromGitHub(accessToken string) model.GithubUser
}

type handlers struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handlers{
		db: db,
	}
}
