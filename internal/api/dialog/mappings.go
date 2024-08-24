package dialog

import "social-network-service/internal/model"

func mapToGetMessagesResponse(messages []*model.Message) GetMessagesResponse {
	items := make([]GetMessagesResponseItem, len(messages))

	for i, msg := range messages {
		item := GetMessagesResponseItem{
			MessageId: int64(msg.MessageId),
			From:      string(msg.FromUserId),
			To:        string(msg.ToUserId),
			Text:      msg.Text,
		}

		items[i] = item
	}

	resp := GetMessagesResponse{
		Messages: items,
	}

	return resp
}
