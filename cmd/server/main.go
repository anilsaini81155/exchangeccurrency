package main

import (
	"log"
	"net/http"

	"github.com/anilsaini81155/exchangeccurrency/handlers"
	"github.com/anilsaini81155/exchangeccurrency/internal/internalredis"
	"github.com/anilsaini81155/exchangeccurrency/internal/websocket"
	"github.com/anilsaini81155/exchangeccurrency/middleware"
)

func main() {
	// Set up Redis client
	redisClient := internalredis.SetupRedis()

	// Set up WebSocket server
	wsServer := websocket.NewWebSocketServer()
	go wsServer.Run()
	go wsServer.RedisSubscriber(redisClient)

	// http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	// http.HandleFunc("/user", handlers.GetUser)
	// router.HandleFunc("/rate", handlers.StoreExchangeRate).Methods("POST").Handler(middleware.JWTAuth(http.HandlerFunc(handlers.StoreExchangeRate)))
	http.Handle("/rate", middleware.JWTAuth(http.HandlerFunc(handlers.ExchangeRates)))
	http.HandleFunc("/ws", wsServer.HandleWebSocket)

	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
