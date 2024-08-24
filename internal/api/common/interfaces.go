package common

import (
	"context"
	"social-network-service/internal/model"

	"github.com/gin-gonic/gin"
)

type AccountService interface {
	RegisterUser(ctx context.Context, cmd *model.RegisterUserCommand) (model.UserId, error)
	Login(ctx context.Context, userId model.UserId, password string) (string, error)
}

type UserService interface {
	GetUser(ctx context.Context, userId model.UserId) (*model.User, error)
	SearchUsers(ctx context.Context, firstName string, secondName string) ([]*model.User, error)
	AddFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId) error
	RemoveFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId) error
}

type PostService interface {
	ReadPosts(ctx context.Context, cmd model.ReadPostsCommand) ([]*model.Post, error)
	CreatePost(ctx context.Context, cmd model.CreatePostCommand) (model.PostId, error)
	GetPost(ctx context.Context, postId model.PostId) (*model.Post, error)
	UpdatePost(ctx context.Context, cmd model.UpdatePostCommand) error
	DeletePost(ctx context.Context, cmd model.DeletePostCommand) error
}

type DialogService interface {
	GetMessages(ctx context.Context, cmd model.GetMessagesCommand) ([]*model.Message, error)
	SendMessage(ctx context.Context, cmd model.SendMessageCommand) error
	NewGetUnreadCountTotal(ctx context.Context, cmd model.NewGetUnreadCountTotalCommand) (int, error)
	NewGetUnreadCount(ctx context.Context, cmd model.NewGetUnreadCountCommand) (int, error)
	MarkMessagesAsRead(ctx context.Context, cmd model.MarkMessagesAsReadCommand) error
}

type JwtService interface {
	GetUserId(c *gin.Context) (model.UserId, error)
}
