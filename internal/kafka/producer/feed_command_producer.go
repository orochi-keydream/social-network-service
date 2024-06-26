package producer

import (
	"encoding/json"
	"social-network-service/internal/kafka/contract"
	"social-network-service/internal/model"

	"github.com/IBM/sarama"
)

type FeedCommandProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewFeedCommandProducer(addrs []string, topic string) (*FeedCommandProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(addrs, cfg)

	if err != nil {
		return nil, err
	}

	p := &FeedCommandProducer{
		topic:    topic,
		producer: producer,
	}

	return p, nil
}

func (p *FeedCommandProducer) PublishAddNewPostToFeedMessage(forUserId model.UserId, postId model.PostId) error {
	payload := contract.AddNewPostToFeedPayload{
		UserId: string(forUserId),
		PostId: string(postId),
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	cmd := contract.UpdateFeedCommand{
		Type:    contract.CommandTypeAddNewPostToFeed,
		Payload: payloadBytes,
	}

	err = p.produce(forUserId, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (p *FeedCommandProducer) PublishUpdatePostInFeedMessage(forUserId model.UserId, postId model.PostId) error {
	payload := contract.UpdatePostInFeedPayload{
		UserId: string(forUserId),
		PostId: string(postId),
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	cmd := contract.UpdateFeedCommand{
		Type:    contract.CommandTypeUpdatePostInFeed,
		Payload: payloadBytes,
	}

	err = p.produce(forUserId, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (p *FeedCommandProducer) PublishRecreateFeedMessage(forUserId model.UserId) error {
	payload := contract.RecreateFeedPayload{
		UserId: string(forUserId),
	}

	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	cmd := contract.UpdateFeedCommand{
		Type:    contract.CommandTypeRecreateFeed,
		Payload: payloadBytes,
	}

	err = p.produce(forUserId, cmd)

	if err != nil {
		return err
	}

	return nil
}

func (p *FeedCommandProducer) produce(userId model.UserId, cmd contract.UpdateFeedCommand) error {
	bytes, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(userId),
		Value: sarama.StringEncoder(bytes),
	}

	_, _, err = p.producer.SendMessage(msg)

	if err != nil {
		return err
	}

	return nil
}
