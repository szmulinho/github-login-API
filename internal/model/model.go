package model

import (
	"os"
)

type GithubUser struct {
	GithubUserID int64  `gorm:"unique_index"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Role         string `json:"role"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
