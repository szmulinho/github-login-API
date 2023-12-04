package model

import (
	"database/sql/driver"
	"encoding/json"
)

type GhUser struct {
	Login       string `json:"login" gorm:"index"`
	AvatarUrl   string `json:"avatar_url" gorm:"index"`
	Role        string `json:"role" gorm:"index"`
}

type PublicRepo struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
}

func (u *GhUser) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *GhUser) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), u)
}
