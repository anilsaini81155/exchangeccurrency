package broadcast

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func publishRateUpdate(client *redis.Client, currencyPair string, rate float64) {
	msg := fmt.Sprintf("%s:%f", currencyPair, rate)
	err := client.Publish(ctx, "rate_updates", msg).Err()
	if err != nil {
		log.Fatalf("Could not publish message: %v", err)
	}
}

func subscribeToRateUpdates(client *redis.Client) {
	pubsub := client.Subscribe(ctx, "rate_updates")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		fmt.Printf("Received update: %s\n", msg.Payload)
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Start a subscriber in a separate goroutine
	go subscribeToRateUpdates(client)

	// Example: Publish a rate update
	publishRateUpdate(client, "USD:EUR", 0.85)
}
