package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the connection to a WebSocket connection
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "Internal server error")

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	// Receive a message
	var msg map[string]string
	err = wsjson.Read(ctx, conn, &msg)
	if err != nil {
		log.Println("Failed to read message:", err)
		return
	}

	fmt.Printf("Received: %v\n", msg)

	// Echo the message back to the client
	err = wsjson.Write(ctx, conn, msg)
	if err != nil {
		log.Println("Failed to write message:", err)
		return
	}

	conn.Close(websocket.StatusNormalClosure, "Goodbye!")
}

func main() {
	http.HandleFunc("/ws", echoHandler)
	fmt.Println("Server is listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
