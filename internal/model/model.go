package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"os"
)

type GitHubLogin struct {
	gorm.Model
	PublicRepos []PublicRepo `gorm:"foreignKey:GitHubLoginID"`
	GithubUser  GithubUser   `gorm:"foreignKey:GitHubLoginID"`
}

type GithubUser struct {
	gorm.Model
	Login       string       `json:"login"`
	Email       string       `json:"email"`
	AvatarUrl   string       `json:"avatarUrl"`
	Followers   int          `json:"followers"`
	AccessToken string       `json:"-"`
	Role        string       `json:"role"`
	PublicRepos []PublicRepo `gorm:"foreignKey:GitHubLoginID"`
}

type PublicRepo struct {
	gorm.Model
	GitHubLoginID   uint `gorm:"foreignKey"`
	GithubUser      GithubUser
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	GithubUserLogin string
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func (u *GithubUser) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *GithubUser) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), u)
}
