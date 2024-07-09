package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-network-service/internal/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AppServiceConfiguration struct {
	TokenGenerator            ITokenGenerator
	UserRepository            IUserRepository
	UserAccountRepository     IUserAccountRepository
	UserFriendRepository      IUserFriendRepository
	DialogRepository          IDialogRepository
	DialogRepositoryTarantool IDialogRepositoryTarantool
	PostRepository            IPostRepository
	FeedCache                 IFeedCache
	FeedCacheNotifier         IFeedCacheNotifier
	PostEventNotifier         IPostEventNotifier
	UserNotifier              IUserNotifier
	TransactionManager        ITransactionManager
}

type AppService struct {
	tokenGenerator            ITokenGenerator
	userRepository            IUserRepository
	userAccountRepository     IUserAccountRepository
	userFriendRepository      IUserFriendRepository
	dialogRepository          IDialogRepository
	dialogRepositoryTarantool IDialogRepositoryTarantool
	postRepository            IPostRepository
	feedCache                 IFeedCache
	cacheNotifier             IFeedCacheNotifier
	postEventNotifier         IPostEventNotifier
	userNotifier              IUserNotifier
	transactionManager        ITransactionManager
}

func NewAppService(config *AppServiceConfiguration) *AppService {
	return &AppService{
		tokenGenerator:            config.TokenGenerator,
		userRepository:            config.UserRepository,
		userAccountRepository:     config.UserAccountRepository,
		userFriendRepository:      config.UserFriendRepository,
		dialogRepository:          config.DialogRepository,
		dialogRepositoryTarantool: config.DialogRepositoryTarantool,
		postRepository:            config.PostRepository,
		feedCache:                 config.FeedCache,
		cacheNotifier:             config.FeedCacheNotifier,
		postEventNotifier:         config.PostEventNotifier,
		userNotifier:              config.UserNotifier,
		transactionManager:        config.TransactionManager,
	}
}

func (s *AppService) GetUser(ctx context.Context, userId model.UserId) (*model.User, error) {
	user, err := s.userRepository.Get(ctx, userId, nil)

	// TODO: Consider using custom error for ErrNoRows.
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.NewNotFoundError(fmt.Sprintf("user %v not found", userId), err)
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *AppService) RegisterUser(ctx context.Context, cu *model.RegisterUserCommand) (model.UserId, error) {
	generatedGuid := uuid.New().String()
	userId := model.UserId(generatedGuid)

	user := &model.User{
		UserId:     userId,
		FirstName:  cu.FirstName,
		SecondName: cu.SecondName,
		Birthdate:  cu.Birthdate,
		Gender:     cu.Gender,
		Biography:  cu.Biography,
		City:       cu.City,
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(cu.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", errors.New("failed to register new user")
	}

	account := &model.UserAccount{
		UserId:       userId,
		PasswordHash: string(hashBytes),
	}

	tx, err := s.transactionManager.Begin(ctx)
	defer s.transactionManager.Rollback(tx)

	if err != nil {
		return "", err
	}

	err = s.userAccountRepository.Add(ctx, account, tx)

	if err != nil {
		return "", err
	}

	err = s.userRepository.Add(ctx, user, tx)

	if err != nil {
		return "", err
	}

	err = s.transactionManager.Commit(tx)

	if err != nil {
		return "", err
	}

	return userId, nil
}

var passwordHashes = map[string][]byte{}

func (s *AppService) RegisterUsers(ctx context.Context, cmds []*model.RegisterUserCommand) error {
	users := make([]*model.User, len(cmds))
	accounts := make([]*model.UserAccount, len(cmds))

	for i, cmd := range cmds {
		generatedGuid := uuid.New().String()
		userId := model.UserId(generatedGuid)

		user := &model.User{
			UserId:     userId,
			FirstName:  cmd.FirstName,
			SecondName: cmd.SecondName,
			Birthdate:  cmd.Birthdate,
			Gender:     cmd.Gender,
			Biography:  cmd.Biography,
			City:       cmd.City,
		}

		fmt.Printf("User %v %v has been created\n", cmd.FirstName, cmd.SecondName)

		hashBytes, found := passwordHashes[cmd.Password]

		if !found {
			hashBytes, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)

			if err != nil {
				return errors.New("failed to register new user")
			}

			passwordHashes[cmd.Password] = hashBytes
		}

		account := &model.UserAccount{
			UserId:       userId,
			PasswordHash: string(hashBytes),
		}

		fmt.Printf("Account %v %v has been created\n", cmd.FirstName, cmd.SecondName)

		users[i] = user
		accounts[i] = account
	}

	tx, err := s.transactionManager.Begin(ctx)
	defer s.transactionManager.Rollback(tx)

	if err != nil {
		return err
	}

	fmt.Println("Saving created accounts")

	err = s.userAccountRepository.AddBulk(ctx, accounts, tx)

	fmt.Println("Accounts have been saved")

	if err != nil {
		return err
	}

	fmt.Println("Saving created users")

	err = s.userRepository.AddBulk(ctx, users, tx)

	if err != nil {
		return err
	}

	fmt.Println("Users have been saved")

	err = s.transactionManager.Commit(tx)

	if err != nil {
		return err
	}

	return nil
}

func (s *AppService) Login(ctx context.Context, userId model.UserId, password string) (string, error) {
	account, err := s.userAccountRepository.Get(ctx, userId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", model.NewNotFoundError(fmt.Sprintf("user %v not found", userId), err)
		default:
			return "", err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))

	// TODO: Consider other possible errors.
	if err != nil {
		return "", model.NewClientError("wrong password", err)
	}

	token, err := s.tokenGenerator.GenerateToken(userId)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AppService) SearchUsers(ctx context.Context, firstName string, secondName string) ([]*model.User, error) {
	users, err := s.userRepository.SearchUsers(ctx, firstName, secondName, nil)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *AppService) SendMessage(ctx context.Context, cmd model.SendMessageCommand) error {
	// TODO: Think about using a transaction.
	_, err := s.userRepository.Get(ctx, cmd.FromUserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", cmd.FromUserId), err)
		default:
			return err
		}
	}

	_, err = s.userRepository.Get(ctx, cmd.ToUserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", cmd.ToUserId), err)
		default:
			return err
		}
	}

	msg := &model.Message{
		FromUserId: cmd.FromUserId,
		ToUserId:   cmd.ToUserId,
		Text:       cmd.Text,
		SentAt:     time.Now().UTC(),
	}

	// _, err = s.dialogRepository.AddMessage(ctx, msg, nil)

	_, err = s.dialogRepositoryTarantool.AddMessage(ctx, msg)

	if err != nil {
		return err
	}

	return nil
}

func (s *AppService) GetMessages(ctx context.Context, cmd model.GetMessagesCommand) ([]*model.Message, error) {
	// messages, err := s.dialogRepository.GetMessages(ctx, cmd.FromUserId, cmd.ToUserId, nil)
	messages, err := s.dialogRepositoryTarantool.GetMessages(ctx, cmd.FromUserId, cmd.ToUserId)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *AppService) ReadPosts(ctx context.Context, cmd model.ReadPostsCommand) ([]*model.Post, error) {
	user, err := s.userRepository.Get(ctx, cmd.UserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.NewNotFoundError(fmt.Sprintf("user %v not found", cmd.UserId), err)
		default:
			return nil, err
		}
	}

	var posts []*model.Post

	if cmd.Offset+cmd.Limit > 1000 {
		posts, err = s.postRepository.GetPostsIncludingFriends(ctx, cmd.UserId, cmd.Offset, cmd.Limit, nil)
	} else {
		posts, err = s.feedCache.GetFeed(user.UserId, cmd.Offset, cmd.Limit)
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *AppService) GetPost(ctx context.Context, postId model.PostId) (*model.Post, error) {
	post, err := s.postRepository.GetPost(ctx, postId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.NewNotFoundError(fmt.Sprintf("post %v not found", postId), err)
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *AppService) CreatePost(ctx context.Context, cmd model.CreatePostCommand) (model.PostId, error) {
	// TODO: Consider using transaction.
	_, err := s.userRepository.Get(ctx, cmd.AuthorUserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.PostId(""), model.NewNotFoundError(fmt.Sprintf("user %v not found", cmd.AuthorUserId), err)
		default:
			return "", err
		}
	}

	postId := model.PostId(uuid.NewString())

	post := &model.Post{
		PostId:       postId,
		PublishedAt:  time.Now().UTC(),
		AuthorUserId: cmd.AuthorUserId,
		Text:         cmd.Text,
	}

	err = s.postRepository.AddPost(ctx, post, nil)

	if err != nil {
		return model.PostId(""), err
	}

	// TODO: Consider using outbox.

	err = s.postEventNotifier.PublishPostCreatedEvent(post)

	if err != nil {
		return model.PostId(""), err
	}

	return post.PostId, nil
}

func (s *AppService) UpdatePost(ctx context.Context, cmd model.UpdatePostCommand) error {
	// TODO: Consider using transaction.
	post, err := s.postRepository.GetPost(ctx, cmd.PostId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("post %v not found", cmd.AuthorUserId), err)
		default:
			return err
		}
	}

	if cmd.AuthorUserId != post.AuthorUserId {
		return model.NewForbiddenError("post does not belong to user", nil)
	}

	post.Text = cmd.Text

	err = s.postRepository.UpdatePost(ctx, post, nil)

	if err != nil {
		return err
	}

	// TODO: Consider using outbox.

	err = s.postEventNotifier.PublishPostUpdatedEvent(post)

	return err
}

func (s *AppService) DeletePost(ctx context.Context, cmd model.DeletePostCommand) error {
	// TODO: Consider using transaction.
	post, err := s.GetPost(ctx, cmd.PostId)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("post %v not found", cmd.AuthorUserId), err)
		default:
			return err
		}
	}

	if cmd.AuthorUserId != post.AuthorUserId {
		return model.NewForbiddenError("post does not beling to user", nil)
	}

	err = s.postRepository.DeletePost(ctx, cmd.PostId, nil)

	if err != nil {
		return err
	}

	// TODO: Consider using outbox.

	err = s.postEventNotifier.PublishPostDeletedEvent(post)

	return err
}

func (s *AppService) AddFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId) error {
	// TODO: Consider using transaction.
	user, err := s.userRepository.Get(ctx, userId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", userId), err)
		default:
			return err
		}
	}

	friendUser, err := s.userRepository.Get(ctx, friendUserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", friendUserId), err)
		default:
			return err
		}
	}

	if user.UserId == friendUser.UserId {
		return model.NewClientError("user IDs must differ", nil)
	}

	err = s.userFriendRepository.AddFriend(ctx, user.UserId, friendUser.UserId, nil)

	if err != nil {
		return err
	}

	// TODO: Consider using outbox.

	err = s.cacheNotifier.PublishRecreateFeedMessage(user.UserId)

	return err
}

func (s *AppService) RemoveFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId) error {
	// TODO: Consider using transaction.
	user, err := s.userRepository.Get(ctx, userId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", userId), err)
		default:
			return err
		}
	}

	friendUser, err := s.userRepository.Get(ctx, friendUserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.NewNotFoundError(fmt.Sprintf("user %v not found", friendUserId), err)
		default:
			return err
		}
	}

	// TODO: What if the row is already removed?
	err = s.userFriendRepository.RemoveFriend(ctx, user.UserId, friendUser.UserId, nil)

	if err != nil {
		return err
	}

	// TODO: Consider using outbox.

	err = s.cacheNotifier.PublishRecreateFeedMessage(user.UserId)

	return err
}

func (s *AppService) SpreadPostCreatedEvent(postId model.PostId, userId model.UserId) error {
	ctx := context.Background()

	friendUserIds, err := s.userFriendRepository.GetSubscribers(ctx, userId, nil)

	if err != nil {
		return err
	}

	s.cacheNotifier.PublishAddNewPostToFeedMessage(userId, postId)

	for _, friendUserId := range friendUserIds {
		err = s.cacheNotifier.PublishAddNewPostToFeedMessage(friendUserId, postId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppService) SpreadPostUpdatedEvent(postId model.PostId, userId model.UserId) error {
	ctx := context.Background()

	friendUserIds, err := s.userFriendRepository.GetSubscribers(ctx, userId, nil)

	if err != nil {
		return err
	}

	s.cacheNotifier.PublishUpdatePostInFeedMessage(userId, postId)

	for _, friendUserId := range friendUserIds {
		err = s.cacheNotifier.PublishUpdatePostInFeedMessage(friendUserId, postId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppService) SpreadPostDeletedEvent(postId model.PostId, userId model.UserId) error {
	ctx := context.Background()

	friendUserIds, err := s.userFriendRepository.GetSubscribers(ctx, userId, nil)

	if err != nil {
		return err
	}

	s.cacheNotifier.PublishRecreateFeedMessage(userId)

	for _, friendUserId := range friendUserIds {
		err = s.cacheNotifier.PublishRecreateFeedMessage(friendUserId)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AppService) AddNewPostToFeedCache(cmd model.AddNewPostToFeedCacheCommand) error {
	ctx := context.Background()

	post, err := s.postRepository.GetPost(ctx, cmd.PostId, nil)

	if err != nil {
		return err
	}

	err = s.feedCache.AddPost(cmd.UserId, post)

	if err != nil {
		return err
	}

	err = s.userNotifier.NotifyNewPostAppeared(ctx, cmd.UserId, post)

	return err
}

func (s *AppService) UpdatePostInFeedCache(cmd model.UpdatePostInFeedCacheCommand) error {
	ctx := context.Background()

	post, err := s.postRepository.GetPost(ctx, cmd.PostId, nil)

	if err != nil {
		return err
	}

	err = s.feedCache.UpdatePost(cmd.UserId, post)

	if err != nil {
		return err
	}

	s.userNotifier.NotifyPostUpdated(ctx, cmd.UserId, post)

	return err
}

func (s *AppService) RecreateFeedCache(cmd model.RecreateFeedCacheCommand) error {
	ctx := context.Background()
	posts, err := s.postRepository.GetPostsIncludingFriends(ctx, cmd.UserId, 0, 1000, nil)

	if err != nil {
		return err
	}

	err = s.feedCache.RecreateFeed(cmd.UserId, posts)

	return err
}
