package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *handlers) Logged(w http.ResponseWriter, r *http.Request, githubData, reposData string) {
	if githubData == "" {
		http.Error(w, "UNAUTHORIZED!", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	responseData := map[string]string{
		"githubData": githubData,
		"reposData":  reposData,
	}

	jsonData, err := json.MarshalIndent(responseData, "", "\t")
	if err != nil {
		log.Println("JSON parse error:", err)
		http.Error(w, "Failed to format JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}
