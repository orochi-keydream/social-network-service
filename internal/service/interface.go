package service

import (
	"context"
	"database/sql"
	"social-network-service/internal/model"
)

type IUserRepository interface {
	Add(ctx context.Context, user *model.User, tx *sql.Tx) error
	AddBulk(ctx context.Context, users []*model.User, tx *sql.Tx) error
	Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.User, error)
	SearchUsers(ctx context.Context, firstName string, secondName string, tx *sql.Tx) ([]*model.User, error)
}

type IUserAccountRepository interface {
	Add(ctx context.Context, account *model.UserAccount, tx *sql.Tx) error
	AddBulk(ctx context.Context, accounts []*model.UserAccount, tx *sql.Tx) error
	Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.UserAccount, error)
}

type IDialogRepository interface {
	AddMessage(ctx context.Context, msg *model.Message, tx *sql.Tx) (model.MessageId, error)
	GetMessages(ctx context.Context, chatId model.ChatId, tx *sql.Tx) ([]*model.Message, error)
}

type IPostRepository interface {
	GetPosts(ctx context.Context, userIds []model.UserId, offset int, limit int, tx *sql.Tx) ([]*model.Post, error)
	GetPostsIncludingFriends(ctx context.Context, userId model.UserId, offset, limit int, tx *sql.Tx) ([]*model.Post, error)
	GetPost(ctx context.Context, postId model.PostId, tx *sql.Tx) (*model.Post, error)
	AddPost(ctx context.Context, post *model.Post, tx *sql.Tx) error
	UpdatePost(ctx context.Context, post *model.Post, tx *sql.Tx) error
	DeletePost(ctx context.Context, postId model.PostId, tx *sql.Tx) error
}

type IUserFriendRepository interface {
	GetFriends(ctx context.Context, userId model.UserId, tx *sql.Tx) ([]model.UserId, error)
	GetSubscribers(ctx context.Context, userId model.UserId, tx *sql.Tx) ([]model.UserId, error)
	AddFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error
	RemoveFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error
}

type ITokenGenerator interface {
	GenerateToken(userId model.UserId) (string, error)
}

type ITransactionManager interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

type IFeedCache interface {
	RecreateFeed(userId model.UserId, posts []*model.Post) error
	GetFeed(userId model.UserId, offset, limit int) ([]*model.Post, error)
	AddPost(userId model.UserId, post *model.Post) error
	UpdatePost(userId model.UserId, post *model.Post) error
}

type IFeedCacheNotifier interface {
	PublishRecreateFeedMessage(forUserId model.UserId) error
	PublishAddNewPostToFeedMessage(forUserId model.UserId, postId model.PostId) error
	PublishUpdatePostInFeedMessage(forUserId model.UserId, postId model.PostId) error
}

type IPostEventNotifier interface {
	PublishPostCreatedEvent(post *model.Post) error
	PublishPostUpdatedEvent(post *model.Post) error
	PublishPostDeletedEvent(post *model.Post) error
}

type IUserNotifier interface {
	NotifyNewPostAppeared(ctx context.Context, userId model.UserId, post *model.Post) error
	NotifyPostUpdated(ctx context.Context, userId model.UserId, post *model.Post) error
}
