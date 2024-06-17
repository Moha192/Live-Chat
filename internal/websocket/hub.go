package ws

import (
	"log"
	"sync"

	"github.com/Moha192/Chat/database"
	"github.com/Moha192/Chat/internal/models"
)

type Hub struct {
	chats      map[int]*chat
	clients    map[int]*client
	broadcast  chan *models.CreateMsgReq
	register   chan *client
	unregister chan *client
	mu         sync.Mutex
}

func NewNub() *Hub {
	return &Hub{
		chats:      make(map[int]*chat),
		clients:    make(map[int]*client),
		broadcast:  make(chan *models.CreateMsgReq),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

type chat struct {
	chatID  int
	clients map[int]*client
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				usersChatIDs, err := database.GetChatIDsByUser(client.userID)
				if err != nil {
					log.Println(err)
					break
				}

				for _, chatID := range usersChatIDs {
					chat, ok := h.chats[chatID]
					if ok {
						delete(chat.clients, client.userID)
					}
				}

				close(client.message)
				delete(h.clients, client.userID)

				log.Printf("logged out, userID: %d", client.userID)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			chat, ok := h.chats[message.ChatID]
			if ok {
				message, err := database.CreateMessage(*message)
				if err != nil {
					log.Println(err)
					break
				}

				for _, client := range chat.clients {
					select {
					case client.message <- &message:
					default:
						close(client.message)
						delete(chat.clients, client.userID)
					}
				}
				log.Printf("Message{content: %s, chatID: %d, userID: %d}", message.Content, message.ChatID, message.UserID)
			}
			h.mu.Unlock()
		}
	}
}
