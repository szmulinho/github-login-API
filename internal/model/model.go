package model

import (
	"os"
)

type GitHubLogin struct {
	PublicRepos []PublicRepo `json:"public_repos"`
	User        struct {
		Login     string `json:"login"`
		Email     string `json:"email"`
		Followers int    `json:"followers"`
	} `json:"user"`
}

type PublicRepo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
