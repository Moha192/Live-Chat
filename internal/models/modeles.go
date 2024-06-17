package models

import "time"

type User struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RespWithUserID struct {
	UserID int `json:"user_id"`
}

type Message struct {
	MessageID int       `json:"message_id"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMsgReq struct {
	ChatID  int    `json:"chat_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}

type CreateDirectChatReq struct {
	UserID       int    `json:"user_id"`
	MemberID     int    `json:"member_id"`
	FirstMessage string `json:"first_message"`
}

type CreateDirectChatResp struct {
	ChatID int `json:"chat_id"`
}

type GetChatsByUserResp struct {
	ChatID      int       `json:"chat_id"`
	ChatName    string    `json:"chat_name"`
	ChatType    string    `json:"chat_type"`
	CreatedAt   time.Time `json:"created_at"`
	LastMessage *Message  `json:"last_message"`
}
