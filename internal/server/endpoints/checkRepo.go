package endpoints

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
)

func (h *handlers) checkRepoAdminAccess(accessToken, repoName string) bool {
	u, err := url.Parse(repoName)
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
