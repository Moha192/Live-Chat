package api

import (
	"github.com/Moha192/Chat/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	CORS(r)

	r.POST("/signup", signUp)
	r.POST("/login", logIn)
	r.GET("/checkjwt", middleware.RequireAuth, check)

	return r
}

func CORS(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Content-Type", "application/json")

		c.Next()
	})
}
