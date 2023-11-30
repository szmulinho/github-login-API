package endpoints

import (
	"net/http"
)

func (h *handlers) HandleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := h.getGithubAccessToken(code)

	reposData := h.getUserData(githubAccessToken, "repos")
	githubData := h.getUserData(githubAccessToken, "")

	h.Logged(w, r, githubData, reposData)
}


