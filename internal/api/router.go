package api

import (
	"github.com/Moha192/Chat/internal/middleware"
	ws "github.com/Moha192/Chat/internal/websocket"
	"github.com/gin-gonic/gin"
)

func SetupRouter(hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	CORS(r)

	r.POST("/signup", signUp)
	r.POST("/login", logIn)
	r.GET("/checkjwt", middleware.RequireAuth, check)

	r.GET("/ws/connect/:user_id", hub.ConnectClient)

	r.GET("/chats/:user_id", hub.GetChatsByUser)
	r.POST("/directChat", hub.CreateDirectChat)
	r.DELETE("/chat/:chat_id", hub.DeleteDirectChat)

	r.GET("/messages/:chat_id", GetMessagesByChat)
	r.PATCH("/messages/:chat_id", SetMessagesStatusToRead)

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
