package producer

import (
	"encoding/json"
	"social-network-service/internal/kafka/contract"
	"social-network-service/internal/model"

	"github.com/IBM/sarama"
)

type PostEventProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewPostEventProducer(addrs []string, topic string) (*PostEventProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(addrs, cfg)

	if err != nil {
		return nil, err
	}

	p := &PostEventProducer{
		topic:    topic,
		producer: producer,
	}

	return p, nil
}

func (p *PostEventProducer) PublishPostCreatedEvent(post *model.Post) error {
	event := contract.PostEvent{
		Type:         contract.EventTypePostCreated,
		PostId:       string(post.PostId),
		AuthorUserId: string(post.AuthorUserId),
	}

	return p.produce(event)
}

func (p *PostEventProducer) PublishPostUpdatedEvent(post *model.Post) error {
	event := contract.PostEvent{
		Type:         contract.EventTypePostUpdated,
		PostId:       string(post.PostId),
		AuthorUserId: string(post.AuthorUserId),
	}

	return p.produce(event)
}

func (p *PostEventProducer) PublishPostDeletedEvent(post *model.Post) error {
	event := contract.PostEvent{
		Type:         contract.EventTypePostDeleted,
		PostId:       string(post.PostId),
		AuthorUserId: string(post.AuthorUserId),
	}

	return p.produce(event)
}

func (p *PostEventProducer) produce(event contract.PostEvent) error {
	bytes, err := json.Marshal(event)

	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(event.AuthorUserId),
		Value: sarama.StringEncoder(bytes),
	}

	_, _, err = p.producer.SendMessage(msg)

	return err
}
