package endpoints

import (
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/gorm"
	"net/http"
)

type Handlers interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	RootHandler(w http.ResponseWriter, r *http.Request)
	GetUserDataHandler(w http.ResponseWriter, r *http.Request)
	getUserFromToken(tokenString string) (*model.GithubUser, error)
}

type handlers struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handlers{
		db: db,
	}
}
