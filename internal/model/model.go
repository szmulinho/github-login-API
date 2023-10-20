package model

import "os"

type GithubUser struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
