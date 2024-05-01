package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"social-network-service/internal/model"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	AddBulk(ctx context.Context, users []*model.User, tx *sql.Tx) error
}

type UserAccountRepository interface {
	AddBulk(ctx context.Context, accounts []*model.UserAccount, tx *sql.Tx) error
}

type TransactionManager interface {
	Begin(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
}

type AdminServiceConfiguration struct {
	UserRepository        UserRepository
	UserAccountRepository UserAccountRepository
	TransactionManager    TransactionManager
}

type AdminService struct {
	userRepository        UserRepository
	userAccountRepository UserAccountRepository
	transactionManager    TransactionManager
}

func NewAdminService(config *AdminServiceConfiguration) *AdminService {
	return &AdminService{
		userRepository:        config.UserRepository,
		userAccountRepository: config.UserAccountRepository,
		transactionManager:    config.TransactionManager,
	}
}

var passwordHashes = map[string][]byte{}

func (s *AdminService) RegisterUsers(ctx context.Context, cmds []*model.RegisterUserCommand) error {
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
