package model

import (
	"os"
)

type GithubUser struct {
	ID           int64        `gorm:"unique_index"`
	Login        string       `json:"login"`
	AvatarUrl    string       `json:"avatar_url"`
	HtmlUrl      string       `json:"html_url"`
	Email        string       `json:"email"`
	Role         string       `json:"role"`
	AccessToken  string       `json:"-"`
	Repositories []Repository `json:"repositories"`
}

type Repository struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
