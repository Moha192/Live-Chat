package api

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func signUp(c *gin.Context) {
	var user models.AuthReq

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if user.Username == "" || len(user.Password) < 4 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := database.CreateUser(user)
	if err != nil {
		if err.Error() == "user already exists" {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusOK)
}

func logIn(c *gin.Context) {
	var user models.AuthReq

	if err := c.Bind(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var (
		err  error
		resp models.RespWithUserID
	)

	resp.UserID, err = database.LogIn(&user)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if resp.UserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "password or email is incorrect",
		})
		return
	}

	tokenString, err := generateJWT(resp.UserID)
	if err != nil {
		log.Println("Error generating JWT:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to generate JWT",
		})
		return
	}

	cookieTime, err := strconv.Atoi(os.Getenv("COOKIE_EXP_TIME"))
	if err != nil {
		log.Println("Error setting cookie time:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to set cookie time",
		})
		return
	}

	c.SetSameSite(http.SameSiteDefaultMode)
	c.SetCookie("Authorization", tokenString, cookieTime, "/", "", false, true)

	c.JSON(http.StatusOK, resp)
}

func check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access": "true",
	})
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
