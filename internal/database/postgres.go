package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() {
	var err error
	DB, err = pgx.Connect(context.Background(), os.Getenv("DB_CONNECTION"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		userid SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL CHECK (LENGTH(username) >= 1),
		password TEXT NOT NULL);`)
	if err != nil {
		log.Fatal(err)
	}
}
