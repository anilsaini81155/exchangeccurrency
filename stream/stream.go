package stream

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func addRateUpdateToStream(client *redis.Client, currencyPair string, rate float64) {
	msg := map[string]interface{}{
		"currencyPair": currencyPair,
		"rate":         rate,
	}
	err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: "rate_stream",
		Values: msg,
	}).Err()
	if err != nil {
		log.Fatalf("Could not add message to stream: %v", err)
	}
}

func consumeRateUpdatesFromStream(client *redis.Client, consumerName string) {
	for {
		streams, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "rate_group",
			Consumer: consumerName,
			Streams:  []string{"rate_stream", ">"},
			Count:    10,
			Block:    0,
		}).Result()
		if err != nil {
			log.Fatalf("Could not read from stream: %v", err)
		}

		for _, stream := range streams {
			for _, message := range stream.Messages {
				currencyPair := message.Values["currencyPair"]
				rate := message.Values["rate"]
				fmt.Printf("Consumer %s received update: %s=%s\n", consumerName, currencyPair, rate)

				// Acknowledge the message
				client.XAck(ctx, "rate_stream", "rate_group", message.ID)
			}
		}
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a consumer group (only need to do this once)
	client.XGroupCreateMkStream(ctx, "rate_stream", "rate_group", "$")

	// Start consumers in separate goroutines
	go consumeRateUpdatesFromStream(client, "consumer-1")
	go consumeRateUpdatesFromStream(client, "consumer-2")

	// Example: Add a rate update to the stream
	addRateUpdateToStream(client, "USD:EUR", 0.85)

	// Prevent the main goroutine from exiting
	select {}
}
