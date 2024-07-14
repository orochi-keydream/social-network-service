package repository

import (
	"context"
	"social-network-service/internal/model"
	"time"

	"github.com/tarantool/go-tarantool/v2"
	"github.com/vmihailenco/msgpack/v5"
)

type getMessageResponse struct {
	messages []getMessageResponseItem
}

type getMessageResponseItem struct {
	messageId  int64
	fromUserId string
	toUserId   string
	sentAt     string
	text       string
}

type addMessageResponse struct {
	messageId int64
}

type DialogRepositoryTarantool struct {
	conn *tarantool.Connection
}

func NewDialogRepositoryTarantool(conn *tarantool.Connection) *DialogRepositoryTarantool {
	return &DialogRepositoryTarantool{
		conn: conn,
	}
}

func (r *DialogRepositoryTarantool) AddMessage(ctx context.Context, msg *model.Message) (model.MessageId, error) {
	args := []interface{}{
		msg.FromUserId,
		msg.ToUserId,
		msg.SentAt.Format(time.RFC3339),
		msg.Text,
	}

	req := tarantool.NewCallRequest("add_message").Args(args)

	var resp addMessageResponse
	err := r.conn.Do(req).GetTyped(&resp)

	if err != nil {
		return model.MessageId(0), err
	}

	return model.MessageId(resp.messageId), nil
}

func (r *DialogRepositoryTarantool) GetMessages(ctx context.Context, curUserId model.UserId, tarUserId model.UserId) ([]*model.Message, error) {
	args := []interface{}{
		curUserId,
		tarUserId,
	}

	req := tarantool.NewCallRequest("get_messages").Args(args)

	var list getMessageResponse
	err := r.conn.Do(req).GetTyped(&list)

	if err != nil {
		return nil, err
	}

	messages := make([]*model.Message, 0, len(list.messages))

	for _, item := range list.messages {
		sentAt, err := time.Parse(time.RFC3339, item.sentAt)

		if err != nil {
			return nil, err
		}

		message := model.Message{
			MessageId:  model.MessageId(item.messageId),
			FromUserId: model.UserId(item.fromUserId),
			ToUserId:   model.UserId(item.toUserId),
			SentAt:     sentAt,
			Text:       item.text,
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

func (m *getMessageResponse) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error

	_, err = d.DecodeArrayLen()

	if err != nil {
		return err
	}

	_, err = d.DecodeMapLen()

	if err != nil {
		return err
	}

	err = d.Skip()

	if err != nil {
		return err
	}

	err = d.Skip()

	if err != nil {
		return err
	}

	err = d.Skip()

	if err != nil {
		return err
	}

	rowCount, err := d.DecodeArrayLen()

	if err != nil {
		return err
	}

	messages := make([]getMessageResponseItem, 0, rowCount)

	for i := 0; i < rowCount; i++ {
		_, err := d.DecodeArrayLen()

		if err != nil {
			return err
		}

		messageId, err := d.DecodeInt64()

		if err != nil {
			return err
		}

		fromUserid, err := d.DecodeString()

		if err != nil {
			return err
		}

		toUserId, err := d.DecodeString()

		if err != nil {
			return err
		}

		sentAt, err := d.DecodeString()

		if err != nil {
			return err
		}

		text, err := d.DecodeString()

		if err != nil {
			return err
		}

		message := getMessageResponseItem{
			messageId:  messageId,
			fromUserId: fromUserid,
			toUserId:   toUserId,
			sentAt:     sentAt,
			text:       text,
		}

		messages = append(messages, message)
	}

	m.messages = messages

	return nil
}

func (m *addMessageResponse) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error

	_, err = d.DecodeArrayLen()

	if err != nil {
		return err
	}

	_, err = d.DecodeArrayLen()

	if err != nil {
		return err
	}

	value, err := d.DecodeInt64()

	if err != nil {
		return err
	}

	m.messageId = value

	return nil
}
