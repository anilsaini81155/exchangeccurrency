package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Connect to the WebSocket server
	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket:", err)
	}
	defer conn.Close(websocket.StatusInternalError, "Internal error")

	// Send a message
	msg := map[string]string{"message": "Hello, Server!"}
	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		log.Fatal("Failed to send message:", err)
	}

	// Receive a message
	var reply map[string]string
	err = wsjson.Read(ctx, conn, &reply)
	if err != nil {
		log.Fatal("Failed to read message:", err)
	}

	fmt.Printf("Received reply: %v\n", reply)

	conn.Close(websocket.StatusNormalClosure, "Goodbye!")
}
