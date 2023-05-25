package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func RunServer() (*server.Server, error) {
	opts := &server.Options{
		JetStream: true,
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}

	go ns.Start()

	if !ns.ReadyForConnections(4 * time.Second) {
		return nil, fmt.Errorf("server is not ready for connection")
	}

	return ns, nil
}

func CreateStream(js nats.JetStreamContext, subj string) error {
	cfg := nats.StreamConfig{
		Name:     subj,
		Subjects: []string{subj},
		Storage:  nats.MemoryStorage,
	}
	_, err := js.AddStream(&cfg)
	if err != nil {
		return err
	}

	return nil
}

func ConnectJetStream(url string) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js, nil
}

func Publisher(js nats.JetStreamContext, subj string) error {
	for {
		id := uuid.New()

		_, err := js.Publish(subj, []byte(id.String()))
		if err != nil {
			return err
		}

		fmt.Printf("next message published %s\n", id.String())

		time.Sleep(getWaitTime())
	}
}

func Subscriber(ctx context.Context, js nats.JetStreamContext, subj string) error {
	sub, err := js.PullSubscribe(subj, "")

	if err != nil {
		return err
	}

	for {
		msgs, err := sub.Fetch(1, nats.Context(ctx))
		if err != nil {
			return err
		}

		fmt.Printf("next message received %s\n", string(msgs[0].Data))
		time.Sleep(getWaitTime())
	}
}

func getWaitTime() time.Duration {
	rand.Seed(time.Now().UnixNano())
	sleepWait := rand.Intn(5000-1000+1) + 1000
	return time.Duration(sleepWait) * time.Millisecond
}
