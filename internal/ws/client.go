package ws

import (
	"log"
	"social-network-service/internal/model"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	userId model.UserId
	send   chan []byte
}

func (c *Client) Read() {
	defer func() {
		c.hub.unregisterClient(c)
	}()

	for {
		_, _, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) Write() {
	for {
		select {
		case msg, ok := <-c.send:
			log.Println("trying to send message")

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
