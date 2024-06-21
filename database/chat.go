package database

import (
	"context"
	"errors"

	"github.com/Moha192/Chat/internal/models"
	"github.com/jackc/pgx/v5"
)

func GetChatsByUser(userID int) ([]models.UsersChatsResp, error) {
	rows, err := DB.Query(context.Background(), `
	SELECT
		c.chat_id,
		CASE
			WHEN c.chat_type = 'direct' THEN (SELECT u.username FROM chat_members cm2 JOIN users u ON cm2.user_id = u.user_id WHERE cm2.chat_id = c.chat_id AND cm2.user_id != $1 LIMIT 1)
			ELSE c.chat_name
		END,
		c.chat_type,
		c.created_at,
		m.message_id,
		m.user_id,
		m.content,
		m.status,
		m.message_type,
		m.created_at as message_created_at
	FROM
		chats c
	JOIN
		chat_members cm ON c.chat_id = cm.chat_id
	JOIN
		messages m ON c.chat_id = m.chat_id
	WHERE
		cm.user_id = $1
		AND m.message_id = (
			SELECT message_id
			FROM messages
			WHERE chat_id = c.chat_id
			ORDER BY created_at DESC
			LIMIT 1
		)
	ORDER BY
		c.created_at DESC;
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []models.UsersChatsResp
	for rows.Next() {
		var chat models.UsersChatsResp
		var msg models.Message

		err := rows.Scan(&chat.ChatID, &chat.ChatName, &chat.ChatType, &chat.CreatedAt, &msg.MessageID, &msg.UserID, &msg.Content, &msg.Status, &msg.Type, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}

		chat.LastMessage = &msg
		chats = append(chats, chat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return chats, nil
}

func GetChatIDsByUser(userID int) ([]int, error) {
	rows, err := DB.Query(context.Background(), "SELECT chat_id FROM chat_members WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	var chatIDs []int
	for rows.Next() {
		var chatID int
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chatIDs, err
}

func CreateDirectChat(userID, memberID int) (int, error) {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	var existingChatID int
	err = tx.QueryRow(context.Background(), `
		SELECT c.chat_id
		FROM chats c
		JOIN chat_members cm1 ON c.chat_id = cm1.chat_id
		JOIN chat_members cm2 ON c.chat_id = cm2.chat_id
		WHERE c.chat_type = $1 AND cm1.user_id = $2 AND cm2.user_id = $3
	`, "direct", userID, memberID).Scan(&existingChatID)

	if err == nil {
		return 0, errors.New("chat already exists")
	}
	if err != pgx.ErrNoRows {
		return 0, err
	}

	var chatID int
	err = tx.QueryRow(context.Background(), "INSERT INTO chats (chat_type) VALUES ($1) RETURNING chat_id", "direct").Scan(&chatID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO chat_members (chat_id, user_id) VALUES($1, $2)", chatID, userID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO chat_members (chat_id, user_id) VALUES($1, $2)", chatID, memberID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func DeleteChat(chatID int) error {
	tx, err := DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "DELETE FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM chat_members WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM chats WHERE chat_id = $1", chatID)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
