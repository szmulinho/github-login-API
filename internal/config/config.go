package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadConfigFromEnv() StorageConfig {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
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
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Dbname, c.Port)
	return connectionString
}
