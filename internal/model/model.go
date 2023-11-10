package model

import (
	"gorm.io/gorm"
	"os"
)

type GitHubLogin struct {
	gorm.Model
	PublicRepos []PublicRepo `gorm:"foreignKey:GitHubLoginID"`
	User        struct {
		Login     string `json:"login"`
		Email     string `json:"email"`
		Followers int    `json:"followers"`
	} `json:"user"`
}

type PublicRepo struct {
	gorm.Model
	GitHubLoginID uint   `gorm:"foreignKey"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
