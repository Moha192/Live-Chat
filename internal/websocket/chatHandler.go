package ws

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Hub) GetChatsByUser(c *gin.Context) {
	stringUserID := c.Param("user_id")
	userID, err := strconv.Atoi(stringUserID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if userID < 1 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, ok := h.clients[userID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not connected",
		})
		return
	}

	chatsResponse, err := database.GetChatsByUser(userID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, chatResp := range chatsResponse {
		if _, ok := h.chats[chatResp.ChatID]; !ok {
			h.chats[chatResp.ChatID] = &chat{
				chatID:  chatResp.ChatID,
				clients: make(map[int]*client),
			}
		}
		h.chats[chatResp.ChatID].clients[userID] = h.clients[userID]
	}

	c.JSON(http.StatusOK, chatsResponse)
}

func (h *Hub) CreateDirectChat(c *gin.Context) {
	var newDirectChatReq models.CreateDirectChatReq
	if err := c.ShouldBindJSON(&newDirectChatReq); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if newDirectChatReq.UserID < 1 || newDirectChatReq.MemberID < 1 || newDirectChatReq.FirstMessage == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	chatID, err := database.CreateDirectChat(newDirectChatReq.UserID, newDirectChatReq.MemberID)
	if err != nil {
		if err.Error() == "chat already exists" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.chats[chatID] = &chat{
		chatID:  chatID,
		clients: make(map[int]*client),
	}

	if client, ok := h.clients[newDirectChatReq.UserID]; ok {
		h.chats[chatID].clients[newDirectChatReq.UserID] = client
	}

	if client, ok := h.clients[newDirectChatReq.MemberID]; ok {
		h.chats[chatID].clients[newDirectChatReq.MemberID] = client
	}

	var msg = models.CreateMsgReq{
		ChatID:  chatID,
		UserID:  newDirectChatReq.UserID,
		Content: newDirectChatReq.FirstMessage,
	}

	log.Printf("Direct chat created:{userID: %d, memberID: %d, chatID: %d}", newDirectChatReq.UserID, newDirectChatReq.MemberID, chatID)
	h.broadcast <- &msg
	c.JSON(http.StatusOK, gin.H{
		"chat_id": chatID,
	})
}

func (h *Hub) DeleteDirectChat(c *gin.Context) {
	strChatID := c.Param("chat_id")
	chatID, err := strconv.Atoi(strChatID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if chatID < 1 {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = database.DeleteChat(chatID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	delete(h.chats, chatID)
	c.Status(http.StatusOK)
}
