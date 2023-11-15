package endpoints

import (
	"github.com/golang-jwt/jwt"
	"github.com/szmulinho/github-login/internal/model"
	"log"
	"net/http"
	"time"
)

func (h *handlers) GenerateToken(w http.ResponseWriter, r *http.Request, githubUserLogin string, isGithubUser bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"githubUserLogin": githubUserLogin,
		"isGithubUser":    isGithubUser,
		"exp":             time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})

	tokenString, err := token.SignedString(model.JwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	} else {
		log.Printf("Token for user %s generated", githubUserLogin)
	}

	return tokenString, nil
}
