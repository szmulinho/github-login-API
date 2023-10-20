package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadFromEnv() StorageConfig {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n")
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	fmt.Printf("Host: %s\n", host)
	fmt.Printf("User: %s\n", user)
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("DBName: %s\n", dbname)
	fmt.Printf("Port: %s\n", port)

	return StorageConfig{
		Host:     host,
		User:     user,
		Password: password,
		Dbname:   dbname,
		Port:     port,
	}
}

type StorageConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
	Port     string `json:"port"`
}

func (c StorageConfig) ConnectionString() string {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		c.Host, c.User, c.Password, c.Dbname, c.Port)
	return connectionString
}
