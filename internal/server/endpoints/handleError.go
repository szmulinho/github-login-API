package endpoints

import (
	"log"
	"net/http"
)

func handleError(w http.ResponseWriter, errMsg string, status int, err error) {
	log.Printf("%s: %v", errMsg, err)
	http.Error(w, errMsg, status)
}
