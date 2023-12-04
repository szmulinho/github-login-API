package endpoints

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/szmulinho/github-login/internal/model"
	"net/http"
	"time"
)

func (h *handlers) GenerateDoctorToken(w http.ResponseWriter, r *http.Request, Login string, isDoctor bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userRole": Login,
		"isDoctor": isDoctor,
		"exp":    time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})

	tokenString, err := token.SignedString(model.JwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	} else {
		fmt.Sprintf("token for user %s generated", model.GhUser{})
	}

	return tokenString, nil
}
