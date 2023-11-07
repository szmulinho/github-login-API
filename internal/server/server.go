package server

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/szmulinho/github-login/internal/server/endpoints"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Run(ctx context.Context, db *gorm.DB) {
	handler := endpoints.NewHandler(db)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/github/login", handler.HandleGitHubLogin)
	router.HandleFunc("/github/callback", handler.HandleCallback)
	http.HandleFunc("/logged", func(w http.ResponseWriter, r *http.Request) {
		endpoints.LoggedHandler(w, r, "")
	})

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://szmul-med.onrender.com"}), // Replace with your React app's origin
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(86400),
	)

	corsRouter := cors(router)

	go func() {
		err := http.ListenAndServe(":8086", corsRouter)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
