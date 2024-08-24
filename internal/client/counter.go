package client

import (
	"context"
	"social-network-service/internal/grpc/counter"
	"social-network-service/internal/model"
)

type CounterClient struct {
	client counter.CounterServiceClient
}

func NewCounterClient(client counter.CounterServiceClient) *CounterClient {
	return &CounterClient{
		client: client,
	}
}

func (c *CounterClient) GetUnreadCountTotal(ctx context.Context, userId model.UserId) (int, error) {
	req := &counter.GetUnreadCountTotalV1Request{
		UserId: string(userId),
	}

	resp, err := c.client.GetUnreadCountTotalV1(ctx, req)

	if err != nil {
		return 0, err
	}

	return int(resp.Count), nil
}

func (c *CounterClient) GetUnreadCount(
	ctx context.Context,
	currentUserId model.UserId,
	chatUserId model.UserId,
) (int, error) {
	req := &counter.GetUnreadCountV1Request{
		CurrentUserId: string(currentUserId),
		ChatUserId:    string(chatUserId),
	}

	resp, err := c.client.GetUnreadCountV1(ctx, req)

	if err != nil {
		return 0, err
	}

	return int(resp.Count), nil
}

func (c *CounterClient) MarkMessagesAsRead(
	ctx context.Context,
	userId model.UserId,
	messageIds []model.MessageId,
) error {
	messageIdList := make([]int64, len(messageIds))

	for i, messageId := range messageIds {
		messageIdList[i] = int64(messageId)
	}

	req := &counter.MarkMessagesAsReadV1Request{
		UserId:     string(userId),
		MessageIds: messageIdList,
	}
	_, err := c.client.MarkMessagesAsReadV1(ctx, req)

	return err
}
