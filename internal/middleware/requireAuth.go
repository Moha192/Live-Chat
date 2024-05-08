package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "no cookie",
		})
		return
	}

	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method: " + token.Header["alg"].(string))
		}

		return []byte("secretkey"), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "JWT expired",
			})
			return
		}

		log.Println(err)
		return
	}
}
