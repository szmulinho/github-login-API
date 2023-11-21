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
	router.HandleFunc("/login", handler.HandleLogin)
	router.HandleFunc("/callback", handler.HandleCallback)
	router.HandleFunc("/register", handler.Register).Methods("POST")
	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		handler.GetUserDataHandler(w, r, tokenString)
	})
	//router.HandleFunc("/user/login", handler.Login).Methods("POST")
	http.HandleFunc("/logged", func(w http.ResponseWriter, r *http.Request) {
		endpoints.Handlers.Logged(handler, w, r, "")

		cors := handlers.CORS(
			handlers.AllowedOrigins([]string{"https://szmul-med.onrender.com", "https://szmul-med.onrender.com/github_user", "https://szmul-med.onrender.com/githubprofile"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type", "Origin", "Accept"}),
			handlers.AllowCredentials(),
			handlers.MaxAge(86400),
		)

		go func() {
			err := http.ListenAndServe(":8086", cors(router))
			if err != nil {
				log.Fatal(err)
			}
		}()

		<-ctx.Done()
	})
}
