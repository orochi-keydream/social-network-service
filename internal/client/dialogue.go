package client

import (
	"context"
	"social-network-service/internal/grpc/dialogue"
	"social-network-service/internal/model"
)

type DialogueClient struct {
	client dialogue.DialogueServiceClient
}

func NewDialogueClient(client dialogue.DialogueServiceClient) *DialogueClient {
	return &DialogueClient{
		client: client,
	}
}

func (c *DialogueClient) SendMessage(ctx context.Context, fromUserId model.UserId, toUserId model.UserId, text string) error {
	req := &dialogue.SendMessageV1Request{
		FromUserId: string(fromUserId),
		ToUserId:   string(toUserId),
		Text:       text,
	}

	_, err := c.client.SendMessageV1(ctx, req)

	return err
}

func (c *DialogueClient) GetMessages(ctx context.Context, fromUserId model.UserId, toUserId model.UserId) ([]*model.Message, error) {
	req := &dialogue.GetMessagesV1Request{
		FromUserId: string(fromUserId),
		ToUserId:   string(toUserId),
	}

	resp, err := c.client.GetMessagesV1(ctx, req)

	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, len(resp.Messages))

	for i, msg := range resp.Messages {
		message := &model.Message{
			MessageId:  model.MessageId(msg.MessageId),
			FromUserId: model.UserId(msg.FromUserId),
			ToUserId:   model.UserId(msg.ToUserId),
			Text:       msg.Text,
		}

		messages[i] = message
	}

	return messages, nil
}
