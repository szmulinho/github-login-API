package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type GithubUser struct {
	Login       string `json:"login" gorm:"index"`
	AvatarUrl   string `json:"avatar_url" gorm:"index"`
	Role        string `json:"role" gorm:"index"`
	AccessToken string `json:"-" gorm:"index"`
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
