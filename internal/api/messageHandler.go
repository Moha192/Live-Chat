package api

import (
	"log"
	"net/http"

	"github.com/Moha192/Chat/database"
	"github.com/gin-gonic/gin"
)

func GetMessagesByChat(c *gin.Context) {
	chatID, err := handleID(c.Param("chat_id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	messages, err := database.GetMessagesByChat(chatID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if messages == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, messages)
}

func SetMessagesStatusToRead(c *gin.Context) {
	messageID, err := handleID(c.Param("message_id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = database.SetMessagesStatusToRead(messageID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func EditMessage(c *gin.Context) {
	messageID, err := handleID(c.Param("message_id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = database.ChangeMessageContent(messageID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteMessage(c *gin.Context) {
	messageID, err := handleID(c.Param("message_id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = database.DeleteMessage(messageID)
	if err != nil {
		if err.Error() == "message not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "message not found",
			})
			return
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}
