package endpoints

import (
	"gorm.io/gorm"
	"net/http"
)

type Handlers interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	Logged(w http.ResponseWriter, r *http.Request, githubData string)
	RegisterDoctor(w http.ResponseWriter, r *http.Request)
	RegisterUser(w http.ResponseWriter, r *http.Request)
	GetUserData(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handlers{
		db: db,
	}
}
