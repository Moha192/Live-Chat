package main

import (
	"context"
	"log"

	"github.com/Moha192/Chat/internal/api"
	"github.com/Moha192/Chat/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	defer database.DB.Close(context.Background())
	log.Println("database connected")

	r := api.SetupRouter()
	log.Println("handlers initialised")
	r.Run(":8080")
}
