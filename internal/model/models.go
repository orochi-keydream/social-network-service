package model

import "time"

type UserId string

type MessageId int64

type PostId string

type Gender int8

const (
	GenderMale Gender = iota + 1
	GenderFemale
)

type User struct {
	UserId     UserId
	FirstName  string
	SecondName string
	Gender     Gender
	Birthdate  time.Time
	Biography  string
	City       string
}

type UserAccount struct {
	UserId       UserId
	PasswordHash string
}

type RegisterUserCommand struct {
	FirstName  string
	SecondName string
	Gender     Gender
	Birthdate  time.Time
	Biography  string
	City       string
	Password   string
}

type Message struct {
	MessageId  MessageId
	SentAt     time.Time
	FromUserId UserId
	ToUserId   UserId
	Text       string
}

type SendMessageCommand struct {
	FromUserId UserId
	ToUserId   UserId
	Text       string
}

type GetMessagesCommand struct {
	FromUserId UserId
	ToUserId   UserId
}

type Post struct {
	PostId       PostId
	PublishedAt  time.Time
	AuthorUserId UserId
	Text         string
}

type ReadPostsCommand struct {
	UserId UserId
	Offset int
	Limit  int
}

type CreatePostCommand struct {
	AuthorUserId UserId
	Text         string
}

type UpdatePostCommand struct {
	AuthorUserId UserId
	PostId       PostId
	Text         string
}

type DeletePostCommand struct {
	AuthorUserId UserId
	PostId       PostId
}
