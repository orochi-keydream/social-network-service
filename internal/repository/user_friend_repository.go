package repository

import (
	"context"
	"database/sql"
	"social-network-service/internal/model"
)

type UserFriendRepository struct {
	cf IConnectionFactory
}

func NewUserFriendRepository(cf IConnectionFactory) *UserFriendRepository {
	return &UserFriendRepository{
		cf: cf,
	}
}

func (r *UserFriendRepository) GetFriends(ctx context.Context, userId model.UserId, tx *sql.Tx) ([]model.UserId, error) {
	const query = "select friend_user_id from user_friends where user_id = $1"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	userIds := []model.UserId{}

	for rows.Next() {
		var friendUserId model.UserId

		err := rows.Scan(&friendUserId)

		if err != nil {
			return nil, err
		}

		userIds = append(userIds, friendUserId)
	}

	return userIds, nil
}

func (r *UserFriendRepository) GetSubscribers(ctx context.Context, userId model.UserId, tx *sql.Tx) ([]model.UserId, error) {
	const query = "select user_id from user_friends where friend_user_id = $1"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, query, userId)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	userIds := []model.UserId{}

	for rows.Next() {
		var userId model.UserId
		err := rows.Scan(&userId)

		if err != nil {
			return nil, err
		}

		userIds = append(userIds, userId)
	}

	return userIds, nil
}

func (r *UserFriendRepository) AddFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error {
	const query = "insert into user_friends (user_id, friend_user_id) values ($1, $2)"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, userId, friendUserId)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserFriendRepository) RemoveFriend(ctx context.Context, userId model.UserId, friendUserId model.UserId, tx *sql.Tx) error {
	const query = "delete from user_friends where user_id = $1 and friend_user_id = $2"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, userId, friendUserId)

	if err != nil {
		return err
	}

	return nil
}
