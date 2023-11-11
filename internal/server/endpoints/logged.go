package endpoints

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func Logged(w http.ResponseWriter, r *http.Request, githubData string) {
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
}
