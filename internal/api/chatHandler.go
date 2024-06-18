package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Moha192/Chat/database"
	chat "github.com/Moha192/Chat/internal/chat"
	"github.com/Moha192/Chat/internal/models"
	"github.com/gin-gonic/gin"
)

func GetChatsByUser(h *chat.Hub, c *gin.Context) {
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

	if _, ok := h.Clients[userID]; !ok {
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
		if _, ok := h.Chats[chatResp.ChatID]; !ok {
			h.Chats[chatResp.ChatID] = &chat.Chat{
				ChatID:  chatResp.ChatID,
				Clients: make(map[int]*chat.Client),
			}
		}
		h.Chats[chatResp.ChatID].Clients[userID] = h.Clients[userID]
	}

	c.JSON(http.StatusOK, chatsResponse)
}

func CreateDirectChat(h *chat.Hub, c *gin.Context) {
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

	h.Chats[chatID] = &chat.Chat{
		ChatID:  chatID,
		Clients: make(map[int]*chat.Client),
	}

	if client, ok := h.Clients[newDirectChatReq.UserID]; ok {
		h.Chats[chatID].Clients[newDirectChatReq.UserID] = client
	}

	if client, ok := h.Clients[newDirectChatReq.MemberID]; ok {
		h.Chats[chatID].Clients[newDirectChatReq.MemberID] = client
	}

	var msg = models.CreateMsgReq{
		ChatID:  chatID,
		UserID:  newDirectChatReq.UserID,
		Content: newDirectChatReq.FirstMessage,
	}

	log.Printf("Direct chat created:{userID: %d, memberID: %d, chatID: %d}", newDirectChatReq.UserID, newDirectChatReq.MemberID, chatID)
	h.Broadcast <- &msg
	c.JSON(http.StatusOK, gin.H{
		"chat_id": chatID,
	})
}

func DeleteDirectChat(h *chat.Hub, c *gin.Context) {
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

	delete(h.Chats, chatID)
	c.Status(http.StatusOK)
}
