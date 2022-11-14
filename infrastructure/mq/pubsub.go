package mq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/aryanugroho/blogapp/internal/logger"
)

type messageQueue struct {
	client *pubsub.Client
	ctx    context.Context
}

const (
	maxAttempt = 5
)

var (
	errFailedPublishMessage = errors.New("failed publish message")
)

func NewMQClient(ctx context.Context, projectID string) (Client, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		logger.Error(ctx, "Error Create Message Queue Client", err)
		return nil, err
	}

	return &messageQueue{
		client: client,
		ctx:    ctx,
	}, nil
}

func (mq *messageQueue) Publish(ctx context.Context, topic string, message []byte) error {

	t := mq.client.Topic(topic)
	result := t.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return err
	}

	if id == "" {
		return errFailedPublishMessage
	}

	return nil
}

func (mq *messageQueue) BulkPublish(topic string, messageList [][]byte) error {

	var results []*pubsub.PublishResult
	var resultErrors []error
	t := mq.client.Topic(topic)
	t.PublishSettings.ByteThreshold = 5000
	t.PublishSettings.CountThreshold = 10
	t.PublishSettings.DelayThreshold = 100 * time.Millisecond

	for _, msg := range messageList {
		result := t.Publish(mq.ctx, &pubsub.Message{
			Data: msg,
		})

		results = append(results, result)
	}

	for _, res := range results {
		id, err := res.Get(mq.ctx)
		if err != nil {
			resultErrors = append(resultErrors, err)
			logger.Error(mq.ctx, fmt.Sprintf("failed to publish message : %s", topic), err)
			continue
		}

		logger.Info(mq.ctx, fmt.Sprintf("success publish message id: %v", id))
	}

	if len(resultErrors) != 0 {
		return fmt.Errorf("process bulk publish not finished: %v", resultErrors[len(resultErrors)-1])
	}

	return nil
}

func (mq *messageQueue) Subscribe(topic, subs string, maxOutstandingMessages int, numGoroutines int, handler Handler) (SubscriptionHandler, error) {
	sub := mq.client.Subscription(subs)
	sub.ReceiveSettings.MaxOutstandingMessages = maxOutstandingMessages
	sub.ReceiveSettings.NumGoroutines = numGoroutines

	err := sub.Receive(mq.ctx, func(ctx context.Context, msg *pubsub.Message) {
		err := handler(topic, subs, &Message{msg})
		if err != nil {
			logger.Error(ctx, fmt.Sprintf("Subscriber Handler Error topic: %s - sub: %s", topic, subs), err)
		}
	})

	if err != nil {
		logger.Error(mq.ctx, fmt.Sprintf("Subscriber Recieve Error topic: %s - sub: %s", topic, subs), err)
		return nil, err
	}

	kch := SubscribeHandler{kc: mq.client, signal: make(chan struct{}, 1)}
	var sc SubscriptionHandler = &kch
	return sc, nil
}

func (mq *messageQueue) Close() error {
	err := mq.client.Close()
	if err != nil {
		return err
	}

	return nil
}

type SubscribeHandler struct {
	kc     *pubsub.Client
	signal chan struct{}
	ctx    context.Context
}

func (sub *SubscribeHandler) Close() error {
	if sub.kc != nil {
		err := sub.kc.Close()
		if err != nil {
			return err
		}
	}

	close(sub.signal)
	return nil
}

func DefaultSubscriberHandlerWrapper(handler func(string, string, []byte) error) Handler {
	return func(topic string, subs string, msg *Message) error {
		err := handler(topic, subs, msg.pubsubMessage.Data)

		var attemp int
		if msg.pubsubMessage.DeliveryAttempt != nil {
			attemp = *msg.pubsubMessage.DeliveryAttempt
		}

		if err != nil && attemp < maxAttempt {
			msg.Nack()
		} else {
			msg.Ack()
		}

		return err
	}
}

func CustomSubscriberHandlerWrapper(handler func(string, string, *Message) error) Handler {
	return func(topic string, subs string, msg *Message) error {
		err := handler(topic, subs, msg)

		if err != nil {
			msg.Nack()
		}

		return err
	}
}

type Message struct {
	pubsubMessage *pubsub.Message
}

func (m *Message) Ack() {
	m.pubsubMessage.Ack()
}

func (m *Message) Nack() {
	m.pubsubMessage.Nack()
}

func (m *Message) Data() []byte {
	return m.pubsubMessage.Data
}
