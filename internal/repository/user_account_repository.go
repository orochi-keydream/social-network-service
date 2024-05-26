package repository

import (
	"context"
	"database/sql"
	"social-network-service/internal/model"

	"github.com/jackc/pgtype"
)

type UserAccountRepository struct {
	cf IConnectionFactory
}

func NewUserAccountRepository(cf IConnectionFactory) *UserAccountRepository {
	return &UserAccountRepository{
		cf: cf,
	}
}

func (r *UserAccountRepository) Add(ctx context.Context, account *model.UserAccount, tx *sql.Tx) error {
	const query = "insert into user_accounts (user_id, password_hash) values ($1, $2)"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	_, err := ec.ExecContext(ctx, query, account.UserId, account.PasswordHash)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserAccountRepository) AddBulk(ctx context.Context, accounts []*model.UserAccount, tx *sql.Tx) error {
	const query = "insert into user_accounts (user_id, password_hash) select * from unnest ($1::text[], $2::text[])"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	count := len(accounts)

	userIds := make([]string, count)
	passwordHashes := make([]string, count)

	for i := 0; i < count; i++ {
		userIds[i] = string(accounts[i].UserId)
		passwordHashes[i] = accounts[i].PasswordHash
	}

	pgUserIds := pgtype.TextArray{}
	pgUserIds.Set(userIds)

	pgPasswordHashes := pgtype.TextArray{}
	pgPasswordHashes.Set(passwordHashes)

	_, err := ec.ExecContext(ctx, query, pgUserIds, pgPasswordHashes)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserAccountRepository) Get(ctx context.Context, userId model.UserId, tx *sql.Tx) (*model.UserAccount, error) {
	const query = "select user_id, password_hash from user_accounts where user_id = $1"

	var ec IExecutionContext

	if tx == nil {
		ec = r.cf.GetMaster()
	} else {
		ec = tx
	}

	row := ec.QueryRowContext(ctx, query, userId)

	err := row.Err()

	if err != nil {
		return nil, err
	}

	dto := struct {
		UserId       string `db:"user_id"`
		PasswordHash string `db:"password_hash"`
	}{}

	err = row.Scan(&dto.UserId, &dto.PasswordHash)

	if err != nil {
		return nil, err
	}

	userAccount := &model.UserAccount{
		UserId:       model.UserId(dto.UserId),
		PasswordHash: dto.PasswordHash,
	}

	return userAccount, nil
}
