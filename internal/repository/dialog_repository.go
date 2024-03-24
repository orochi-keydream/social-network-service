package repository

import (
	"context"
	"database/sql"
	"social-network-service/internal/model"
)

type DialogRepository struct {
	db *sql.DB
}

func NewDialogRepository(db *sql.DB) *DialogRepository {
	return &DialogRepository{
		db: db,
	}
}

func (r *DialogRepository) AddMessage(ctx context.Context, msg *model.Message, tx *sql.Tx) (model.MessageId, error) {
	const query = "insert into messages (sent_at, from_user_id, to_user_id, text) values ($1, $2, $3, $4) returning message_id"

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	row := ec.QueryRowContext(ctx, query, msg.SentAt, msg.FromUserId, msg.ToUserId, msg.Text)

	if row.Err() != nil {
		return 0, row.Err()
	}

	var messageId model.MessageId
	err := row.Scan(&messageId)

	if err != nil {
		return 0, err
	}

	return messageId, nil
}

func (r *DialogRepository) GetMessages(ctx context.Context, fromUserId model.UserId, toUserId model.UserId, tx *sql.Tx) ([]*model.Message, error) {
	const query = "select message_id, sent_at, from_user_id, to_user_id, text from messages where (from_user_id = $1 and to_user_id = $2) or (from_user_id = $2 and to_user_id = $1) order by sent_at desc"

	var ec ExecutionContext

	if tx == nil {
		ec = r.db
	} else {
		ec = tx
	}

	rows, err := ec.QueryContext(ctx, query, fromUserId, toUserId)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	messages := []*model.Message{}

	for rows.Next() {
		var msg model.Message
		err = rows.Scan(&msg.MessageId, &msg.SentAt, &msg.FromUserId, &msg.ToUserId, &msg.Text)

		if err != nil {
			return nil, err
		}

		messages = append(messages, &msg)
	}

	return messages, nil
}
