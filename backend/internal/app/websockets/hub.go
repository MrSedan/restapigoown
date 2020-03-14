package websockets

import (
	"time"

	"github.com/MrSedan/restapigoown/backend/internal/app/model"
)

//Hub is a room for websockets
type Hub struct {
	ID         string
	server     *Server
	clients    map[*Client]bool
	broadcast  chan *model.Message
	register   chan *Client
	unregister chan *Client
}

//NewHub creating room for websockets
func NewHub(id string, serv *Server) *Hub {
	return &Hub{
		ID:         id,
		server:     serv,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan *model.Message, 100),
		register:   make(chan *Client, 100),
		unregister: make(chan *Client, 100),
	}
}

//Run running websocket room
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if _, ok := h.clients[client]; !ok {
				h.clients[client] = true
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				if len(h.clients) == 0 && len(h.register) == 0 {
					h.server.RemHub <- h
					return
				}
			}
		case message := <-h.broadcast:
			ti := time.Now().Unix()
			h.server.db.User().NewMessage(message.FromID, message.ToID, message.Body, ti)
			message.Time = ti
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
		}

	}
}
