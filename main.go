package main

import (
	"fmt"
	"github.com/szmulinho/common/utils"
	"github.com/szmulinho/github-login/internal/database"
	"github.com/szmulinho/github-login/internal/server"
	"log"
)

func main() {
	fmt.Println("Starting the application...")

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("connetcting to database: %v", err)
	}

	ctx, _, wait := utils.Gracefully()

	server.Run(ctx ,db)

	wait()
}
