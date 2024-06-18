package chat

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

type Client struct {
	UserID     int
	Connection *websocket.Conn
	Message    chan *models.Message
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

	Client := &Client{
		UserID:     userID,
		Connection: conn,
		Message:    make(chan *models.Message),
	}

	h.Register <- Client
	go Client.readPump(h)
	go Client.writePump()
}

func (c *Client) readPump(hub *Hub) {
	defer func() {
		c.Connection.Close()
		hub.Unregister <- c
	}()

	for {
		_, messageWS, err := c.Connection.ReadMessage()
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

		if chat, ok := hub.Chats[message.ChatID]; ok {
			if _, ok = chat.Clients[message.UserID]; !ok {
				c.wsRespondWithError("Client not connetcted")
				continue
			}
		} else {
			c.wsRespondWithError("Chat not found")
			continue
		}

		hub.Broadcast <- &message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.Connection.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		c.Connection.WriteJSON(message)
	}
}

func (c *Client) wsRespondWithError(err string) {
	errorMsg := struct {
		Error string `json:"error"`
	}{
		Error: err,
	}

	c.Connection.WriteJSON(errorMsg)
}
