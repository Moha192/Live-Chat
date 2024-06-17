package database

import (
	"context"

	"github.com/Moha192/Chat/internal/models"
)

func CreateMessage(chatReq models.CreateMsgReq) (models.Message, error) {
	var msg = models.Message{
		ChatID:  chatReq.ChatID,
		UserID:  chatReq.UserID,
		Content: chatReq.Content,
		Type:    "user",
		Status:  "delivered",
	}

	err := DB.QueryRow(context.Background(), "INSERT INTO messages (chat_id, user_id, content) VALUES ($1, $2, $3) RETURNING message_id, created_at", chatReq.ChatID, chatReq.UserID, chatReq.Content).Scan(&msg.MessageID, &msg.CreatedAt)
	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

func GetMessagesByChat(chatID int) ([]models.Message, error) {
	rows, err := DB.Query(context.Background(), "SELECT message_id, chat_id, user_id, content, status, created_at FROM messages WHERE chat_id = $1 ORDER BY created_at ASC", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.MessageID, &msg.ChatID, &msg.UserID, &msg.Content, &msg.Status, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func SetMessagesStatusToRead(messageID int) error {
	_, err := DB.Exec(context.Background(), `
		UPDATE messages
		SET status = 'read'
		WHERE status = 'delivered'
			AND created_at <= (SELECT created_at FROM messages WHERE message_id = $1)
	`, messageID)
	if err != nil {
		return err
	}

	return nil
}
