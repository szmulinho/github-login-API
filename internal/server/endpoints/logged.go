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
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	var prettyJSON bytes.Buffer
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}