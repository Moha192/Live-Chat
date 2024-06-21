package main

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/api"
	hub "github.com/Moha192/Chat/internal/hub"
)

func main() {
	time.Sleep(time.Second * 5) // wait for docker database connection

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.InitDB()
	defer database.DB.Close(context.Background())
	log.Println("database connected")

	hub := hub.NewHub()
	go hub.Run()

	r := api.SetupRouter(hub)
	r.Run(":8080")
}
