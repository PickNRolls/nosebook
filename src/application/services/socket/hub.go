package socket

import (
	"log"

	"github.com/google/uuid"
)

type Hub struct {
	clients map[uuid.UUID]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: map[uuid.UUID]*Client{},
	}
}

func (this *Hub) UserClient(userId uuid.UUID) *Client {
	return this.clients[userId]
}

func (this *Hub) Subscribe(client *Client) {
	this.clients[client.userId] = client
	log.Printf("Subscribe new client for user(id:%v)\n", client.userId)
}

func (this *Hub) Unsubscribe(userId uuid.UUID) {
	if _, has := this.clients[userId]; has {
		log.Printf("Unsubscribe client for user(id:%v)\n", userId)
		client := this.clients[userId]
		delete(this.clients, userId)
		close(client.Send())
	}
}

func (this *Hub) Broadcast(message []byte) {
	for userId, client := range this.clients {
		select {
		case client.Send() <- message:

		default:
			this.Unsubscribe(userId)
		}
	}
}
