package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Moha192/Chat/database"
	"github.com/gin-gonic/gin"
)

func GetMessagesByChat(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chat_id"))
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
	chatID, err := strconv.Atoi(c.Param("chat_id"))
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

	err = database.SetMessagesStatusToRead(chatID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusOK)
}

func DeleteMessage(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("message_id"))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if messageID < 1 {
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
