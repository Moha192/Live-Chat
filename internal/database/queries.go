package database

import (
	"context"
	"errors"

	"github.com/Moha192/Chat/internal/models"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(user *models.User) error {
	var existingUserID int
	err := DB.QueryRow(context.Background(), "SELECT userid FROM users WHERE username = $1", user.Username).Scan(&existingUserID)
	if err == nil {
		return errors.New("user already exists")
	}

	if err != pgx.ErrNoRows {
		return err
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	err = DB.QueryRow(context.Background(), "INSERT INTO users(username, password) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING userid", user.Username, user.Password).Scan(&user.UserID)
	if err != nil {
		return err
	}

	return nil
}

func LogIn(user *models.User) (bool, error) {
	var hashedPassword []byte
	err := DB.QueryRow(context.Background(), "SELECT userid, password FROM users WHERE username = $1", user.Username).Scan(&user.UserID, &hashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.Password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
