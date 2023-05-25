package internal

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsConfig struct {
	Url          string
	Subject      string
	Stream       string
	StartSeq     uint64
	IsLeaf       bool
	Ctx          context.Context
	Consumer     string
	Conn         *nats.Conn
	JSCtx        nats.JetStreamContext
	Subscription *nats.Subscription
}

// ConnectToNATS makes a new nats connection
func (nc *NatsConfig) ConnectToNATS() error {
	var err error

	// connect options
	opts := []nats.Option{
		nats.Name("NATS Sample Subscriber"),
	}

	// connect to NATS
	nc.Conn, err = nats.Connect(nc.Url, opts...)
	if err != nil {
		return err
	}

	var apiPrefix string
	if nc.IsLeaf {
		apiPrefix = "$JS.main.API"
	}

	// obtain JetStreamContext
	nc.JSCtx, err = nc.Conn.JetStream(nats.APIPrefix(apiPrefix))
	if err != nil {
		return err
	}

	return nil
}

// SubscribeToSubject subscribes to the provided subject
func (nc *NatsConfig) SubscribeToSubject() error {
	var err error

	opts := []nats.SubOpt{}

	// if the start sequence id provided, then set it
	if nc.StartSeq > 0 {
		opts = append(opts, nats.StartSequence(nc.StartSeq))
	}

	// if both start sequence id and consumer provider, then try to drop existing consumer
	if nc.StartSeq > 0 && len(nc.Consumer) > 0 {
		err = nc.JSCtx.DeleteConsumer(nc.Stream, nc.Consumer)
		if err != nil && err != nats.ErrConsumerNotFound {
			return err
		}
	}

	// subscribe the stream
	sub, err := nc.JSCtx.PullSubscribe(nc.Subject, nc.Consumer, opts...)
	if err != nil {
		return err
	}

	// fetch infinitely
	for {
		msgs, err := sub.Fetch(10, nats.Context(nc.Ctx))
		if err != nil {
			return err
		}

		// process message pack
		for _, msg := range msgs {
			meta, err := msg.Metadata()
			fmt.Printf("%d: %s\n", meta.Sequence.Stream, string(msg.Data))
			if err = msg.Ack(); err != nil {
				return err
			}
		}
	}
}
