package ws

import (
	"context"
	"log"
	"social-network-service/internal/model"
	"sync"
)

type Hub struct {
	mux     sync.Mutex
	clients map[model.UserId]*clientList
	send    chan message
}

type clientList struct {
	mux     sync.RWMutex
	clients []*Client
}

type message struct {
	userId model.UserId
	msg    []byte
}

func NewHub() *Hub {
	return &Hub{
		mux:     sync.Mutex{},
		clients: make(map[model.UserId]*clientList),
		send:    make(chan message),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case msg, ok := <-h.send:
			if !ok {
				log.Println("hub: Run: channel closed")
				return
			}

			clientList, found := h.clients[msg.userId]

			if !found {
				log.Println("No clients for sending")
				continue
			}

			log.Printf("Got client %v for sending", msg.userId)

			clients := clientList.GetClients()

			for _, client := range clients {
				client.send <- msg.msg
			}
		case <-ctx.Done():
			close(h.send)
		}
	}
}

func (h *Hub) PushMessage(userId model.UserId, msg []byte) error {
	wsMsg := message{
		userId: userId,
		msg:    msg,
	}

	h.send <- wsMsg

	return nil
}

func (h *Hub) registerClient(client *Client) {
	log.Printf("Registering client %v", client.userId)

	h.mux.Lock()
	clients, found := h.clients[client.userId]

	if !found {
		clientList := &clientList{}
		h.clients[client.userId] = clientList
		clients = clientList
	}

	h.mux.Unlock()

	clients.Add(client)
}

func (h *Hub) unregisterClient(client *Client) {
	log.Printf("Unregistering client %v", client.userId)

	h.mux.Lock()
	defer h.mux.Unlock()

	clientList, found := h.clients[client.userId]

	if !found {
		return
	}

	clientList.Remove(client)

	if len(clientList.clients) == 0 {
		delete(h.clients, client.userId)
	}
}

func (l *clientList) Add(client *Client) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.clients = append(l.clients, client)
}

func (l *clientList) Remove(client *Client) {
	l.mux.Lock()
	defer l.mux.Unlock()

	var foundIdx = -1

	for i, curClient := range l.clients {
		if curClient != client {
			continue
		}

		foundIdx = i
	}

	if foundIdx == -1 {
		return
	}

	copy(l.clients[foundIdx:], l.clients[foundIdx+1:])
	l.clients[len(l.clients)-1] = nil
	l.clients = l.clients[:len(l.clients)-1]
}

func (l *clientList) GetClients() []*Client {
	l.mux.RLock()
	defer l.mux.RUnlock()

	clients := make([]*Client, len(l.clients))
	copy(clients, l.clients)

	return clients
}
