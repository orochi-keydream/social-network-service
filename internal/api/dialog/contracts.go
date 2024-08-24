package dialog

type GetMessagesResponse struct {
	// List of messages.
	Messages []GetMessagesResponseItem `json:"messages"`
}

type GetMessagesResponseItem struct {
	// ID of the message.
	MessageId int64 `json:"messageId"`

	// ID of the message sender in UUIDv4 format.
	From string `json:"from"`

	// ID of the message recipient in UUIDv4 format.
	To string `json:"to"`

	// Content of the message.
	Text string `json:"text"`
}

type SendMessageRequest struct {
	// Content of the message.
	Text string `json:"text" binding:"required"`
}

type NewGetUnreadCountTotalResponse struct {
	Count int `json:"count"`
}

type NewGetUnreadCountResponse struct {
	Count int `json:"count"`
}

type MarkMessagesAsReadRequest struct {
	MessageIds []int64 `json:"messageIds"`
}
