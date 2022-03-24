package main

import (
	"context"
	"log"
	"time"

	"github.com/hamba/avro"
	"github.com/twmb/franz-go/pkg/kgo"
)

func InsertOrderConsumer() {
	opts := []kgo.Opt{
		kgo.SeedBrokers(getRedPandaHosts()...),
		kgo.ConsumeTopics(
			topicOrdersInsertAVRO,
		),
	}
	redPandaClient, err := kgo.NewClient(opts...)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	startTime := time.Now()
	schema := avro.MustParse(schemaOrder)

	var counter int64

consumerLoop:
	for {
		fetches := redPandaClient.PollRecords(context.Background(), 1000000)
		for _, fetchErr := range fetches.Errors() {
			log.Printf("error consuming from topic: topic=%s, partition=%d, err=%v\n",
				fetchErr.Topic, fetchErr.Partition, fetchErr.Err)
			break consumerLoop
		}

		records := fetches.Records()
		for _, record := range records {
			counter++
			var order orderAVRO
			err := avro.Unmarshal(schema, record.Value, &order)
			if err != nil {
				log.Println(err)
				continue
			}

			if counter == 1000000 {
				log.Printf("%d %+v", counter, order)
				log.Println("Speed:", time.Since(startTime).Nanoseconds())
			}
		}
	}
}
