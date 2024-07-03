package ws

import (
	"net/http"
	"social-network-service/internal/model"
	"social-network-service/internal/service"
)

func HandleHub(hub *Hub, jwtService *service.JwtService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := jwtService.CheckAccessFromRequest(r)

		if err != nil {
			_, matches := err.(*model.UnauthenticatedError)

			if matches {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = Handle(hub, userId, w, r)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func Handle(hub *Hub, userId model.UserId, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	client := &Client{
		hub:    hub,
		conn:   conn,
		userId: userId,
		send:   make(chan []byte),
	}

	hub.registerClient(client)

	go client.Read()
	go client.Write()

	return nil
}
