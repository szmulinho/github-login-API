package config

import (
	"fmt"
	"os"
)

func LoadFromEnv() StorageConfig {
	host := os.Getenv("DB_H")
	user := os.Getenv("DB_U")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_N")
	port := os.Getenv("DB_P")

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
