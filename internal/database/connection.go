package database

import (
	"github.com/szmulinho/github-login/internal/config"
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	conn := config.LoadFromEnv()
	connectionString := conn.ConnectionString()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&model.Repository{}); err != nil {
		return nil, err
	}

	return db, nil
}
