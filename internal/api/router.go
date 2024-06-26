package api

import (
	hub "github.com/Moha192/Chat/internal/hub"
	"github.com/Moha192/Chat/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hub *hub.Hub) *gin.Engine {
	r := gin.Default()

	CORS(r)

	r.POST("/signup", SignUp)
	r.POST("/login", LogIn)
	r.GET("/checkjwt", middleware.RequireAuth, check)

	r.GET("/ws/connect/:user_id", hub.ConnectClient)

	r.GET("/chats/:user_id", func(c *gin.Context) {
		GetChatsByUser(hub, c)
	})
	r.POST("/directChat", func(c *gin.Context) {
		CreateDirectChat(hub, c)
	})
	r.DELETE("/chat/:chat_id", func(c *gin.Context) {
		DeleteDirectChat(hub, c)
	})

	r.GET("/messages/:chat_id", GetMessagesByChat)
	r.PATCH("/messages/status/:message_id", SetMessagesStatusToRead)
	r.PATCH("/message/:message_id", EditMessage)
	r.DELETE("/message/:message_id", DeleteMessage)

	return r
}

func CORS(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	})
}
