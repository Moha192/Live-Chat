package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/Moha192/Chat/internal/models"
)

func CreateUser(user models.AuthReq) error {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}

	var userID int
	err = DB.QueryRow(context.Background(), "INSERT INTO users(username, password_hash) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING user_id", user.Username, hashedPassword).Scan(&userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("user already exists")
		}
		return err
	}

	return nil
}

func LogIn(user *models.AuthReq) (int, error) {
	var (
		userID       int
		passwordHash []byte
	)
	err := DB.QueryRow(context.Background(), "SELECT user_id, password_hash FROM users WHERE username = $1", user.Username).Scan(&userID, &passwordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(user.Password))
	if err != nil {
		return 0, nil
	}
	return userID, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func UserExists(userID int) (bool, error) {
	var exists bool
	err := DB.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
