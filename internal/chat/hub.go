package chat

import (
	"log"
	"sync"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/models"
)

type Hub struct {
	Chats      map[int]*Chat
	Clients    map[int]*Client
	Broadcast  chan *models.CreateMsgReq
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
}

func NewNub() *Hub {
	return &Hub{
		Chats:      make(map[int]*Chat),
		Clients:    make(map[int]*Client),
		Broadcast:  make(chan *models.CreateMsgReq),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

type Chat struct {
	ChatID  int
	Clients map[int]*Client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mu.Lock()
			h.Clients[client.UserID] = client
			h.Mu.Unlock()

		case client := <-h.Unregister:
			h.Mu.Lock()
			if _, ok := h.Clients[client.UserID]; ok {
				usersChatIDs, err := database.GetChatIDsByUser(client.UserID)
				if err != nil {
					log.Println(err)
					break
				}

				for _, chatID := range usersChatIDs {
					chat, ok := h.Chats[chatID]
					if ok {
						delete(chat.Clients, client.UserID)
					}
				}

				close(client.Message)
				delete(h.Clients, client.UserID)

				log.Printf("logged out, userID: %d", client.UserID)
			}
			h.Mu.Unlock()

		case message := <-h.Broadcast:
			h.Mu.Lock()
			chat, ok := h.Chats[message.ChatID]
			if ok {
				message, err := database.CreateMessage(*message)
				if err != nil {
					log.Println(err)
					break
				}

				for _, client := range chat.Clients {
					select {
					case client.Message <- &message:
					default:
						close(client.Message)
						delete(chat.Clients, client.UserID)
					}
				}
				log.Printf("Message{content: %s, chatID: %d, userID: %d}", message.Content, message.ChatID, message.UserID)
			}
			h.Mu.Unlock()
		}
	}
}
