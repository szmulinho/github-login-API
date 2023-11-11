package endpoints

import (
	"golang.org/x/oauth2"
	"net/http"
	"os"
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

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := oauthConfig2.AuthCodeURL("", oauth2.AccessTypeOnline)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
