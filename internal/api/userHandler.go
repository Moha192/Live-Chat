package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Moha192/Chat/internal/database"
	"github.com/Moha192/Chat/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func signUp(c *gin.Context) {
	var user models.User

	if c.Bind(&user) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	} else if user.Username == "" || len(user.Password) < 4 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := database.SignUp(&user)
	if err != nil {
		if err.Error() == "user already exists" {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, user)
}

func logIn(c *gin.Context) {
	var user models.User

	if c.Bind(&user) != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return

	} else if user.Username == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	findUser, err := database.LogIn(&user)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !findUser {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "password or email is incorrect",
		})
		return
	}

	tokenString, err := generateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to generate JWT",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, user)
}

func check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access": "true",
	})
}

func generateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Second * 20).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
