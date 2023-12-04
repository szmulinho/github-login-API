package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if err := json.Indent(&prettyJSON, []byte(githubData), "", "\t"); err != nil {
		log.Printf("JSON parse error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err := fmt.Fprintf(w, string(prettyJSON.Bytes()))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successful response for request from %s", r.RemoteAddr)
}
