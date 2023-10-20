package model

import (
	"os"
)

type Drug struct {
	DrugID int64  `json:"drug_id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Price  int64  `json:"price"`
}

var Drugs []Drug

type Exception struct {
	Message string `json:"message"`
}

var JwtKey = []byte(os.Getenv("JWT_KEY"))
