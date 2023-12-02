package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type GithubUser struct {
	gorm.Model
	Email       string `json:"email"`
	Login       string `json:"login"`
	Name	    string  `json:"name"`
	AvatarUrl   string `json:"avatar_url"`
	Followers   int    `json:"followers"`
	Role        string `json:"role"`
	AccessToken string `json:"-"`
}

type PublicRepo struct {
	gorm.Model
	GitHubLoginID uint   `gorm:"foreignKey" json:"gitHubLoginID"`
	ID            uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string `json:"name"`
	Description   string `json:"description"`
}

func (u *GithubUser) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *GithubUser) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), u)
}
