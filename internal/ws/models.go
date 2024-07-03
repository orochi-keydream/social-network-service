package ws

import "encoding/json"

type RouterMessage struct {
	UserId  string          `json:"userId"`
	Type    PayloadType     `json:"type"`
	Message json.RawMessage `json:"message"`
}

type PayloadType string

const (
	PayloadTypeNewPostAppearedPayload PayloadType = "NewPostAppeared"
	PayloadTypePostUpdatedPayload     PayloadType = "PostUpdated"
)

type NewPostAppearedPayload struct {
	AuthorUserId string `json:"authorUserId"`
	PostId       string `json:"postId"`
	Text         string `json:"text"`
	PublishedAt  string `json:"publishedAt"`
}

type PostUpdatedPayload struct {
	PostId string `json:"postId"`
	Text   string `json:"text"`
}

type PostDeletedPayload struct {
	PostId string `json:"postId"`
}

type HubMessage struct {
	Type    PayloadType     `json:"type"`
	Message json.RawMessage `json:"message"`
}
