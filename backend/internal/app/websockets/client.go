package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type message struct {
	ID   int    `json:"id"`
	Body string `json:"msg"`
}

//Client ...
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan *message
	doneCh chan bool
}

//ServeWs serving websocket connection
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan *message, 256), doneCh: make(chan bool)}
	client.hub.register <- client

	go client.listenRead()
	go client.listenWrite()
}

func (c *Client) listenRead() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case <-c.doneCh:
			c.doneCh <- true
			return
		default:
			msg := message{}
			err := c.conn.ReadJSON(&msg)
			if err != nil {
				c.doneCh <- true
				return
			}
			c.hub.broadcast <- &msg
		}
	}
}

func (c *Client) listenWrite() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		select {
		case <-c.doneCh:
			return
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte("Closed conn"))
				return
			}

			c.conn.WriteJSON(message)
		}
	}
}
