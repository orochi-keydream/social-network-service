package dialog

import "social-network-service/internal/model"

func mapToGetMessagesResponse(messages []*model.Message) GetMessagesResponse {
	items := []GetMessagesResponseItem{}

	for _, msg := range messages {
		item := GetMessagesResponseItem{
			From: string(msg.FromUserId),
			To:   string(msg.ToUserId),
			Text: msg.Text,
		}

		items = append(items, item)
	}

	resp := GetMessagesResponse{
		Messages: items,
	}

	return resp
}
