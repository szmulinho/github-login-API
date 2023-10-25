package endpoints

import (
	"context"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"net/http"
)

type Handlers interface {
	exchangeGitHubCodeForToken(ctx context.Context, code string) (*oauth2.Token, error)
	getGitHubUserInfo(ctx context.Context, token *oauth2.Token) (*model.GithubUser, error)
	GitHubLoginHandler(w http.ResponseWriter, r *http.Request)
	GitHubCallbackHandler(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handlers{
		db: db,
	}
}
