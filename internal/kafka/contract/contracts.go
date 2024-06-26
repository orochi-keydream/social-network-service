package contract

import "encoding/json"

type UpdateFeedCommand struct {
	Type    CommandType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type AddNewPostToFeedPayload struct {
	UserId string `json:"userId"`
	PostId string `json:"postId"`
}

type UpdatePostInFeedPayload struct {
	UserId string `json:"userId"`
	PostId string `json:"postId"`
}

type RecreateFeedPayload struct {
	UserId string `json:"userId"`
}

type CommandType string

const (
	CommandTypeAddNewPostToFeed CommandType = "AddNewPostToFeed"
	CommandTypeUpdatePostInFeed CommandType = "UpdatePostInFeed"
	CommandTypeRecreateFeed     CommandType = "RecreateFeed"
)

type PostEvent struct {
	Type         EventType `json:"type"`
	PostId       string    `json:"postId"`
	AuthorUserId string    `json:"authorUserId"`
}

type EventType string

const (
	EventTypePostCreated EventType = "EventTypePostCreated"
	EventTypePostUpdated EventType = "EventTypePostUpdated"
	EventTypePostDeleted EventType = "EventTypePostDeleted"
)
