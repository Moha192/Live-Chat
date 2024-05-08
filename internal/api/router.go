package api

import (
	"github.com/Moha192/Chat/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/signUp", signUp)
	r.POST("/logIn", logIn)
	r.GET("/check", middleware.RequireAuth, check)

	return r
}
