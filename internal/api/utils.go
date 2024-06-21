package api

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func handleID(strID string) (int, error) {
	ID, err := strconv.Atoi(strID)
	if err != nil {
		return 0, err
	}

	if ID < 1 {
		return 0, errors.New("ID < 1")
	}

	return ID, nil
}

func generateJWT(userID int) (string, error) {
	JWTExpTime, err := strconv.Atoi(os.Getenv("JWT_EXP_TIME"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Second * time.Duration(JWTExpTime)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
