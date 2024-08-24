package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"social-network-service/internal/kafka/contract"
	"social-network-service/internal/model"
	"social-network-service/internal/service"
	"sync"

	"github.com/IBM/sarama"
)

func UsePostEventConsumer(ctx context.Context, addrs []string, topic string, c *PostEventConsumer, wg *sync.WaitGroup) error {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	cg, err := sarama.NewConsumerGroup(addrs, "social-network-service", cfg)

	if err != nil {
		return err
	}

	topics := []string{topic}

	go func() {
		defer wg.Done()

		for {
			err = cg.Consume(ctx, topics, c)

			if err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}

				log.Panicln("Something wrong with the consumer")
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func NewPostEventConsumer(as *service.AppService) *PostEventConsumer {
	return &PostEventConsumer{
		appService: as,
	}
}

type PostEventConsumer struct {
	appService *service.AppService
}

func (p *PostEventConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (p *PostEventConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (p *PostEventConsumer) ConsumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-cgc.Messages():
			if !ok {
				log.Println("Message channel was closed")
				return nil
			}

			log.Printf("Handling message with offset %v from %v topic\n", msg.Offset, msg.Topic)

			err := p.handle(msg.Value)

			if err != nil {
				log.Println(err)
				continue
			}

			cgs.MarkMessage(msg, "")
		case <-cgs.Context().Done():
			log.Println("ConsumeClaim: cancellation requested")
			return nil
		}
	}
}

func (p *PostEventConsumer) handle(msg []byte) error {
	var postEvent contract.PostEvent
	err := json.Unmarshal(msg, &postEvent)

	if err != nil {
		return err
	}

	switch postEvent.Type {
	case contract.EventTypePostCreated:
		err = p.appService.SpreadPostCreatedEvent(model.PostId(postEvent.PostId), model.UserId(postEvent.AuthorUserId))
	case contract.EventTypePostUpdated:
		err = p.appService.SpreadPostUpdatedEvent(model.PostId(postEvent.PostId), model.UserId(postEvent.AuthorUserId))
	case contract.EventTypePostDeleted:
		err = p.appService.SpreadPostDeletedEvent(model.PostId(postEvent.PostId), model.UserId(postEvent.AuthorUserId))
	}

	return err
}
