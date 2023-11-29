package config

import (
	"fmt"
	"os"
)

func LoadFromEnv() StorageConfig {
	DB_HOST := os.Getenv("DB_HOST")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PORT := os.Getenv("DB_PORT")

	fmt.Printf("DB_HOST: %s\n", DB_HOST)
	fmt.Printf("DB_USER: %s\n", DB_USER)
	fmt.Printf("DB_PASSWORD: %s\n", DB_PASSWORD)
	fmt.Printf("DB_NAME: %s\n", DB_NAME)
	fmt.Printf("DB_PORT: %s\n", DB_PORT)

	return StorageConfig{
		DB_HOST:     DB_HOST,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_NAME:     DB_NAME,
		DB_PORT:     DB_PORT,
	}
}

type StorageConfig struct {
	DB_HOST     string `json:"DB_HOST"`
	DB_USER     string `json:"DB_USER"`
	DB_PASSWORD string `json:"DB_PASSWORD"`
	DB_NAME     string `json:"DB_NAME"`
	DB_PORT     string `json:"DB_PORT"`
}

func (c StorageConfig) ConnectionString() string {
	connectionString := fmt.Sprintf("DB_HOST=%s\n DB_USER=%s\n DB_PASSWORD=%s\n DB_NAME=%s\n DB_PORT=%s\n",
		c.DB_HOST, c.DB_USER, c.DB_PASSWORD, c.DB_NAME, c.DB_PORT)
	return connectionString
}
