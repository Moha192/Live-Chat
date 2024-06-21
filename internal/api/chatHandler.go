package api

import (
	"log"
	"net/http"

	"github.com/Moha192/Chat/database"
	hub "github.com/Moha192/Chat/internal/hub"
	"github.com/Moha192/Chat/internal/models"
	"github.com/gin-gonic/gin"
)

func GetChatsByUser(h *hub.Hub, c *gin.Context) {
	userID, err := handleID(c.Param("user_id"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	h.Mu.Lock()
	client, userExists := h.Clients[userID]
	h.Mu.Unlock()

	if !userExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not connected",
		})
		return
	}

	usersChats, err := database.GetChatsByUser(userID)
	if err != nil {
		log.Println("Error fetching chats for user:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//add chats to the hub if they are not there
	h.Mu.Lock()
	for _, chat := range usersChats {
		if _, ok := h.Chats[chat.ChatID]; !ok {
			h.Chats[chat.ChatID] = &hub.Chat{
				ChatID:  chat.ChatID,
				Clients: make(map[int]*hub.Client),
			}
		}
		//add client to chat
		h.Chats[chat.ChatID].Clients[userID] = client
	}
	h.Mu.Unlock()

	c.JSON(http.StatusOK, usersChats)
}

func CreateDirectChat(h *hub.Hub, c *gin.Context) {
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

	//add chat to the hub
	h.Chats[chatID] = &hub.Chat{
		ChatID:  chatID,
		Clients: make(map[int]*hub.Client),
	}

	//add members of chat
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

func DeleteDirectChat(h *hub.Hub, c *gin.Context) {
	chatID, err := handleID(c.Param("chat_id"))
	if err != nil {
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
