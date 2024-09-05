package socket

import "log"

type Client interface {
	Send() chan []byte
}

type hub struct {
	clients map[Client]struct{}
}

func NewHub() *hub {
	return &hub{
		clients: map[Client]struct{}{},
	}
}

func (this *hub) Subscribe(client Client) {
	this.clients[client] = struct{}{}
	log.Printf("New hub client: %v\n", client)
}

func (this *hub) Unsubscribe(client Client) {
	if _, has := this.clients[client]; has {
		log.Printf("Unsubscribe client: %v\n", client)
		delete(this.clients, client)
		close(client.Send())
	}
}

func (this *hub) Broadcast(message []byte) {
	for client := range this.clients {
		select {
		case client.Send() <- message:

		default:
			this.Unsubscribe(client)
		}
	}
}
