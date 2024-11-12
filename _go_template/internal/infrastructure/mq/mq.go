package mq

import (
	"context"
	"io"
)

type SubscriptionHandler interface {
	io.Closer
}

type Client interface {
	io.Closer
	Publish(ctx context.Context, topicName string, data []byte) error
	BulkPublish(string, [][]byte) error
	// Subscribe receive handler function that will be passed pubsub message to parameter.
	//
	// IMPORTANT: You always have to ack() or nack() the message to acknowledge pubsub message
	// for alternative you can use DefaultSubscriberHandlerWrapper that receive function which return error and treat it as acknowledgement
	Subscribe(
		topic string,
		subscription string,
		maxOutstandingMessages int,
		numGoroutines int,
		handler Handler,
	) (SubscriptionHandler, error)
}

type Handler = func(string, string, *Message) error
