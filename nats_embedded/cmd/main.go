package main

import (
	"context"
	"sync"
	"time"

	"embedded/internal"
)

const subject = "mymessages"

func main() {
	ns, err := internal.RunServer()
	if err != nil {
		panic(err)
	}
	defer ns.Shutdown()

	wg := new(sync.WaitGroup)

	// connect subscriber to the jetstream
	subNCJS, err := internal.ConnectJetStream(ns.ClientURL())
	if err != nil {
		panic(err)
	}

	// make new stream
	err = internal.CreateStream(subNCJS, subject)
	if err != nil {
		panic(err)
	}

	// run subscriber goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
		defer cancel()
		err = internal.Subscriber(ctx, subNCJS, subject)
		if err != nil {
			panic(err)
		}
	}()

	// connect publisher to the jetstream
	pubNCJS, err := internal.ConnectJetStream(ns.ClientURL())
	if err != nil {
		panic(err)
	}

	// run publisher goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = internal.Publisher(pubNCJS, subject)
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
