package ws

import (
	"context"
	"encoding/json"
	"log"
	"social-network-service/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserNotifier struct {
	client *redis.Client
	hub    *Hub
}

func NewUserNotifier(client *redis.Client, hub *Hub) *UserNotifier {
	return &UserNotifier{
		client: client,
		hub:    hub,
	}
}

func (r *UserNotifier) NotifyNewPostAppeared(ctx context.Context, userId model.UserId, post *model.Post) error {
	payload := NewPostAppearedPayload{
		AuthorUserId: string(post.AuthorUserId),
		PostId:       string(post.PostId),
		PublishedAt:  post.PublishedAt.Format(time.RFC3339),
		Text:         post.Text,
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return r.notify(ctx, userId, PayloadTypeNewPostAppearedPayload, payloadBytes)
}

func (r *UserNotifier) NotifyPostUpdated(ctx context.Context, userId model.UserId, post *model.Post) error {
	payload := PostUpdatedPayload{
		PostId: string(post.PostId),
		Text:   post.Text,
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return r.notify(ctx, userId, PayloadTypePostUpdatedPayload, payloadBytes)
}

func (r *UserNotifier) Subscribe(ctx context.Context) {
	pubsub := r.client.Subscribe(ctx, "hub")

	for msg := range pubsub.Channel() {
		var message RouterMessage
		err := json.Unmarshal([]byte(msg.Payload), &message)

		if err != nil {
			log.Println(err)
			continue
		}

		hubMessage := HubMessage{
			Type:    message.Type,
			Message: message.Message,
		}

		hubMessageBytes, err := json.Marshal(hubMessage)

		if err != nil {
			log.Println(err)
			continue
		}

		err = r.hub.PushMessage(model.UserId(message.UserId), hubMessageBytes)

		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func (r *UserNotifier) notify(ctx context.Context, userId model.UserId, payloadType PayloadType, payloadBytes []byte) error {
	message := RouterMessage{
		UserId:  string(userId),
		Type:    payloadType,
		Message: payloadBytes,
	}

	messageBytes, err := json.Marshal(message)

	if err != nil {
		return err
	}

	cmd := r.client.Publish(ctx, "hub", messageBytes)

	return cmd.Err()
}
