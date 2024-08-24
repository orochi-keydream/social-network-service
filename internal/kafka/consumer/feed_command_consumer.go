package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"social-network-service/internal/kafka/contract"
	"social-network-service/internal/model"
	"social-network-service/internal/service"
	"sync"

	"github.com/IBM/sarama"
)

func UseFeedCommandConsumer(ctx context.Context, addrs []string, topic string, c *FeedCommandConsumer, wg *sync.WaitGroup) error {
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

type FeedCommandConsumer struct {
	appService *service.AppService
}

func NewFeedCommandConsumer(as *service.AppService) *FeedCommandConsumer {
	return &FeedCommandConsumer{
		appService: as,
	}
}

func (f *FeedCommandConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (f *FeedCommandConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (f *FeedCommandConsumer) ConsumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-cgc.Messages():
			if !ok {
				log.Println("Message channel was closed")
				return nil
			}

			log.Printf("Handling message with offset %v from %v topic\n", msg.Offset, msg.Topic)

			err := f.handle(msg.Value)

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

func (f *FeedCommandConsumer) handle(msg []byte) error {
	var cmd contract.UpdateFeedCommand
	err := json.Unmarshal(msg, &cmd)

	if err != nil {
		return err
	}

	switch cmd.Type {
	case contract.CommandTypeAddNewPostToFeed:
		var payload contract.AddNewPostToFeedPayload
		err := json.Unmarshal(cmd.Payload, &payload)

		if err != nil {
			return err
		}

		cmd := model.AddNewPostToFeedCacheCommand{
			UserId: model.UserId(payload.UserId),
			PostId: model.PostId(payload.PostId),
		}

		return f.appService.AddNewPostToFeedCache(cmd)
	case contract.CommandTypeUpdatePostInFeed:
		var payload contract.UpdatePostInFeedPayload
		err := json.Unmarshal(cmd.Payload, &payload)

		if err != nil {
			return err
		}

		cmd := model.UpdatePostInFeedCacheCommand{
			UserId: model.UserId(payload.UserId),
			PostId: model.PostId(payload.PostId),
		}

		return f.appService.UpdatePostInFeedCache(cmd)
	case contract.CommandTypeRecreateFeed:
		var payload contract.RecreateFeedPayload
		err := json.Unmarshal(cmd.Payload, &payload)

		if err != nil {
			return err
		}

		cmd := model.RecreateFeedCacheCommand{
			UserId: model.UserId(payload.UserId),
		}

		return f.appService.RecreateFeedCache(cmd)
	default:
		// TODO: Check if an actual type is used in this string.
		return fmt.Errorf("no handler for type %v", reflect.TypeOf(cmd.Payload))
	}
}
