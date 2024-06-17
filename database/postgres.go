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

	_, err = DB.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users (
		user_id SERIAL PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS chats (
		chat_id SERIAL PRIMARY KEY,
		chat_name VARCHAR(50),
		chat_type varchar(10) NOT NULL CHECK (chat_type IN ('direct', 'group', 'channel')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS chat_members (
		chat_id INT NOT NULL,
		user_id INT NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (chat_id, user_id),
		FOREIGN KEY (chat_id) REFERENCES chats(chat_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id));`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS messages (
		message_id SERIAL PRIMARY KEY,
		chat_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		message_type varchar(10) NOT NULL  DEFAULT 'user' CHECK (message_type IN ('user', 'system')),
		status VARCHAR(10) NOT NULL DEFAULT 'delivered' CHECK (status IN ('delivered', 'read')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (chat_id) REFERENCES chats(chat_id),
		FOREIGN KEY (user_id) REFERENCES users(user_id));`)
	if err != nil {
		log.Fatal(err)
	}
}
