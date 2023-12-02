package database

import (
	"github.com/szmulinho/github-login/internal/config"
	"github.com/szmulinho/github-login/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.PublicRepo{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&model.GithubUser{}); err != nil {
		return err
	}
	return nil
}

func Connect() (*gorm.DB, error) {
	conn := config.LoadFromEnv()
	connectionString := conn.ConnectionString()

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err := Migrate(db); err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
