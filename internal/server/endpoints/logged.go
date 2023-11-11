package endpoints

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if err := json.Indent(&prettyJSON, []byte(githubData), "", "\t"); err != nil {
		log.Println("JSON parse error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}
