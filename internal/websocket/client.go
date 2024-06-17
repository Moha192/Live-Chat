package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type client struct {
	userID     int
	connection *websocket.Conn
	message    chan *models.Message
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Hub) ConnectClient(c *gin.Context) {
	strUserID := c.Param("user_id")
	userID, err := strconv.Atoi(strUserID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userExists, err := database.UserExists(userID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "uesr not exists",
		})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	client := &client{
		userID:     userID,
		connection: conn,
		message:    make(chan *models.Message),
	}

	h.register <- client
	go client.readPump(h)
	go client.writePump()
}

func (c *client) readPump(hub *Hub) {
	defer func() {
		c.connection.Close()
		hub.unregister <- c
	}()

	for {
		_, messageWS, err := c.connection.ReadMessage()
		if err != nil {
			break
		}

		var message models.CreateMsgReq

		err = json.Unmarshal(messageWS, &message)
		if err != nil {
			log.Println(err)
			break
		}

		if message.ChatID < 1 || message.UserID < 1 || message.Content == "" {
			c.wsRespondWithError("Bad request")
			continue
		}

		if chat, ok := hub.chats[message.ChatID]; ok {
			if _, ok = chat.clients[message.UserID]; !ok {
				c.wsRespondWithError("Client not connetcted")
				continue
			}
		} else {
			c.wsRespondWithError("Chat not found")
			continue
		}

		hub.broadcast <- &message
	}
}

func (c *client) writePump() {
	defer func() {
		c.connection.Close()
	}()

	for {
		message, ok := <-c.message
		if !ok {
			c.connection.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.connection.WriteJSON(message)
	}
}

func (c *client) wsRespondWithError(err string) {
	errorMsg := struct {
		Error string `json:"error"`
	}{
		Error: err,
	}

	c.connection.WriteJSON(errorMsg)
}
