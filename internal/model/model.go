package model

import (
	"os"
)

type GithubUser struct {
	ID       int64
	Username string
	Email    string
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
