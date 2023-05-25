package main

import (
	"context"
	"flag"
	"log"
	"time"

	"subscriber/internal"
)

func parseFlags() *internal.NatsConfig {
	var config internal.NatsConfig

	flag.Uint64Var(&config.StartSeq, "start-seq", 0, "start sequence number")
	flag.StringVar(&config.Url, "url", "nats://localhost:4222", "nats URL")
	flag.StringVar(&config.Subject, "subject", "testsubject", "messages subject")
	flag.StringVar(&config.Stream, "stream", "teststream", "nats stream")
	flag.BoolVar(&config.IsLeaf, "leaf", false, "is leaf connection")
	flag.StringVar(&config.Consumer, "consumer", "", "durable consumer name")
	flag.Parse()

	return &config
}

func main() {
	config := parseFlags()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1*time.Hour))
	defer cancel()

	config.Ctx = ctx

	err := config.ConnectToNATS()
	if err != nil {
		log.Fatal(err)
	}

	err = config.SubscribeToSubject()
	if err != nil {
		log.Fatal(err)
	}
}
