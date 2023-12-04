package endpoints

import (
	"golang.org/x/oauth2"
	"net/http"
)

var Config = oauth2.Config{
	ClientID:     "ea8a922aeb25cedebae5",
	ClientSecret: "e16a220ee7376a0e41c1d934251ae09af0e1f787",
	Scopes:       []string{"public_repo", "read:user", "user:email", "user:follow"},
	RedirectURL:  "https://szmul-med.onrender.com",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	},
}

func (h *handlers) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirectURL := Config.AuthCodeURL("", oauth2.AccessTypeOffline)
	http.Redirect(w, r, redirectURL, http.StatusFound)
}
