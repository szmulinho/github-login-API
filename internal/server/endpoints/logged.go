package endpoints

import (
	"bytes"
	"encoding/json"
	"github.com/szmulinho/github-login/internal/model"
	"log"
	"net/http"
)

func (h *handlers) Logged(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if err != nil {
		log.Println("JSON parse error:", err)
		http.Error(w, "Failed to format JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(prettyJSON.Bytes())

	githubUser := model.GithubUser{}
	tokenString, err := h.GenerateToken(w, r, githubUser, true)
	if err != nil {
		log.Println("Error generating token:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.getUserFromToken(tokenString)

}
