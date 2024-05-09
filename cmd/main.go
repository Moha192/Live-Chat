package main

import (
	"context"
	"log"
	"time"

	"github.com/Moha192/Chat/internal/api"
	"github.com/Moha192/Chat/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	time.Sleep(time.Second * 1) // wait for docker database connection

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	defer database.DB.Close(context.Background())
	log.Println("database connected")

	r := api.SetupRouter()
	r.Run(":8080")
}
