package server

import (
	"context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/szmulinho/github-login/internal/server/endpoints"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func Run(ctx context.Context, db *gorm.DB) {
	handler := endpoints.NewHandler(db)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/login", handler.HandleLogin)
	router.HandleFunc("/callback", handler.HandleCallback)
	router.HandleFunc("/register", handler.RegisterUser).Methods("POST")
	router.HandleFunc("/register_doctor", handler.RegisterDoctor).Methods("POST")
	router.HandleFunc("/user", handler.GetUserData).Methods("GET")

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8086" // Default port if not provided
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"https://szmul-med.onrender.com", "https://szmul-med.onrender.com/github_user", "https://szmul-med.onrender.com/githubprofile"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type", "Origin", "Accept"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(86400),
	)

	go func() {
		err := http.ListenAndServe(":"+port, cors(router))
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
}
