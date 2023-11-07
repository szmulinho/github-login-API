package model

import (
	"os"
)

type GithubUser struct {
	ID           int64         `gorm:"unique_index"`
	Login        string        `json:"username"`
	AvatarUrl    string        `json:"avatar_url"`
	HtmlUrl      string        `json:"html_url"`
	Email        string        `json:"email"`
	Role         string        `json:"role"`
	AccessToken  string        `json:"-"`
	Repositories []*Repository `gorm:"many2many:user_repositories;"`
}

type Repository struct {
	ID           int64  `gorm:"primaryKey;autoIncrement"`
	Name         string `json:"name"`
	Permission   string `json:"permission"`
	GithubUserID int64  `gorm:"many2many:user_repositories;"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
