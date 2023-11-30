package endpoints

import (
	"golang.org/x/oauth2"
	"net/http"
)

var oauthConfig = oauth2.Config{
	ClientID:     "33f5f8298ded51f76f30",
	ClientSecret: "6598a705034fbe74a6b0b3a4e08af79a7dfa3eac",
	Scopes:       []string{"public_repo", "read:user", "user:email", "user:follow"},
	RedirectURL:  "https://szmul-med.onrender.com/github_register",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := oauthConfig.AuthCodeURL("", oauth2.AccessTypeOnline)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
