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
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/drug", handler.HandleLogin).Methods("POST")
	router.HandleFunc("/drug", handler.HandleCallback).Methods("POST")
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Authorization", "Content-Type"}),
		handlers.ExposedHeaders([]string{}),
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
}
