package endpoints

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/szmulinho/github-login/internal/model"
	"golang.org/x/oauth2"
)

func (h *handlers) checkRepoAdminAccess(accessToken string, user model.GithubUser) bool {
	owner := user.Login
	repoName := "szmul-med"

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})))

	_, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		return false
	}

	return true
}
