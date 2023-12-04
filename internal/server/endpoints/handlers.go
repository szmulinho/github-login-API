package endpoints

import (
	"gorm.io/gorm"
	"net/http"
)

type Handlers interface {
	HandleLogin(w http.ResponseWriter, r *http.Request)
	HandleCallback(w http.ResponseWriter, r *http.Request)
	GenerateUserToken(w http.ResponseWriter, r *http.Request, Login string, isUser bool) (string, error)
	GenerateDoctorToken(w http.ResponseWriter, r *http.Request, Login string, isDoctor bool) (string, error)
}

type handlers struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handlers {
	return &handlers{
		db: db,
	}
}
