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
	Response    Response     `gorm:"foreignKey:GitHubLoginID"`
}

type Response struct {
	gorm.Model
	Login       string `json:"login"`
	Email       string `json:"email"`
	AvatarUrl   string `json:"avatar_url"`
	Followers   int    `json:"followers"`
	Role        string `json:"role"`
	AccessToken string `json:"-"`
}

type GithubUser struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Role      string `json:"role"`
}

type User struct {
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Role      string `json:"role"`
	Followers int    `json:"followers"`
	Email     string `json:"email"`
}

type PublicRepo struct {
	gorm.Model
	GitHubLoginID uint   `gorm:"foreignKey"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func (u *Response) Value() (driver.Value, error) {
	return json.Marshal(u)
}

type LoginResponse struct {
	GithubUser GithubUser `json:"githubUser"`
	Token      string     `json:"token"`
}

func (u *Response) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), u)
}
