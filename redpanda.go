package main

import (
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

func getRedPandaHosts() []string {
	return []string{
		"redpanda:9092",
	}
}

func getRedPandaClient() *kgo.Client {
	opts := []kgo.Opt{
		kgo.SeedBrokers(getRedPandaHosts()...),
		kgo.ConsumeTopics(
			topicOrdersInsertAVRO,
			topicOrdersUpsertAVRO,
			topicTradesInsertAVRO,
		),
	}
	redPandaClient, err := kgo.NewClient(opts...)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	return redPandaClient
}
