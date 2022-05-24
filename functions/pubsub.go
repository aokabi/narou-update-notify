package functions

import (
	"context"
	"log"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}

func NotifyPubSub(ctx context.Context, m PubSubMessage) error {
	log.Printf("Hello pubsub")
	return nil
}
