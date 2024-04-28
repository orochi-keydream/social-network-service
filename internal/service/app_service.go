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

type UserRepository interface {
	Add(ctx context.Context, user *model.User, tx *sql.Tx) error
	AddBulk(ctx context.Context, users []*model.User, tx *sql.Tx) error
	Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.User, error)
	SearchUsers(ctx context.Context, firstName string, secondName string, tx *sql.Tx) ([]*model.User, error)
}

type UserAccountRepository interface {
	Add(ctx context.Context, account *model.UserAccount, tx *sql.Tx) error
	AddBulk(ctx context.Context, accounts []*model.UserAccount, tx *sql.Tx) error
	Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.UserAccount, error)
}

type DialogRepository interface {
	AddMessage(ctx context.Context, msg *model.Message, tx *sql.Tx) (model.MessageId, error)
	GetMessages(ctx context.Context, fromUserId model.UserId, toUserId model.UserId, tx *sql.Tx) ([]*model.Message, error)
}

type PostRepository interface {
	GetPosts(ctx context.Context, userIds []model.UserId, offset int, limit int, tx *sql.Tx) ([]*model.Post, error)
	GetPost(ctx context.Context, postId model.PostId, tx *sql.Tx) (*model.Post, error)
	AddPost(ctx context.Context, post *model.Post, tx *sql.Tx) error
	UpdatePost(ctx context.Context, post *model.Post, tx *sql.Tx) error
	DeletePost(ctx context.Context, postId model.PostId, tx *sql.Tx) error
}

type UserFriendRepository interface {
	GetFriends(ctx context.Context, userId model.UserId, tx *sql.Tx) ([]model.UserId, error)
	AddFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error
	RemoveFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error
}

type TokenGenerator interface {
	GenerateToken(userId model.UserId) (string, error)
}

type TransactionManager interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

type AppServiceConfiguration struct {
	TokenGenerator        TokenGenerator
	UserRepository        UserRepository
	UserAccountRepository UserAccountRepository
	UserFriendRepository  UserFriendRepository
	DialogRepository      DialogRepository
	PostRepository        PostRepository
	TransactionManager    TransactionManager
}

type AppService struct {
	tokenGenerator        TokenGenerator
	userRepository        UserRepository
	userAccountRepository UserAccountRepository
	userFriendRepository  UserFriendRepository
	dialogRepository      DialogRepository
	postRepository        PostRepository
	transactionManager    TransactionManager
}

func NewAppService(config *AppServiceConfiguration) *AppService {
	return &AppService{
		tokenGenerator:        config.TokenGenerator,
		userRepository:        config.UserRepository,
		userAccountRepository: config.UserAccountRepository,
		userFriendRepository:  config.UserFriendRepository,
		dialogRepository:      config.DialogRepository,
		postRepository:        config.PostRepository,
		transactionManager:    config.TransactionManager,
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

	_, err = s.dialogRepository.AddMessage(ctx, msg, nil)

	if err != nil {
		return err
	}

	return nil
}

func (s *AppService) GetMessages(ctx context.Context, cmd model.GetMessagesCommand) ([]*model.Message, error) {
	messages, err := s.dialogRepository.GetMessages(ctx, cmd.FromUserId, cmd.ToUserId, nil)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *AppService) ReadPosts(ctx context.Context, cmd model.ReadPostsCommand) ([]*model.Post, error) {
	// TODO: Consider using transactions.
	user, err := s.userRepository.Get(ctx, cmd.UserId, nil)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, model.NewNotFoundError(fmt.Sprintf("user %v not found", cmd.UserId), err)
		default:
			return nil, err
		}
	}

	friendUserIds, err := s.userFriendRepository.GetFriends(ctx, user.UserId, nil)

	userIds := append([]model.UserId{user.UserId}, friendUserIds...)

	if err != nil {
		return nil, err
	}

	posts, err := s.postRepository.GetPosts(ctx, userIds, cmd.Offset, cmd.Limit, nil)

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

	post := model.Post{
		PostId:       postId,
		PublishedAt:  time.Now().UTC(),
		AuthorUserId: cmd.AuthorUserId,
		Text:         cmd.Text,
	}

	s.postRepository.AddPost(ctx, &post, nil)

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

	s.postRepository.UpdatePost(ctx, post, nil)

	return nil
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

	return nil
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

	// TODO: What if the row is already present?
	err = s.userFriendRepository.AddFriend(ctx, user.UserId, friendUser.UserId, nil)

	if err != nil {
		return err
	}

	return nil
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

	return nil
}
