package model

import (
	"os"
)

type GithubUser struct {
	ID          int64  `gorm:"unique_index"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	AccessToken string `json:"-"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
