package websocket

import (
	"context"
	"log"
	"net/http"
	"sync"

	// "github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type WebSocketServer struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan string
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Mutex      sync.Mutex
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan string),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (s *WebSocketServer) Run() {
	for {
		select {
		case conn := <-s.Register:
			s.Mutex.Lock()
			s.Clients[conn] = true
			s.Mutex.Unlock()

		case conn := <-s.Unregister:
			s.Mutex.Lock()
			if _, ok := s.Clients[conn]; ok {
				delete(s.Clients, conn)
				conn.Close(websocket.StatusNormalClosure, "Closed")
			}
			s.Mutex.Unlock()

		case message := <-s.Broadcast:
			s.Mutex.Lock()
			for conn := range s.Clients {
				err := wsjson.Write(context.Background(), conn, message)
				if err != nil {
					conn.Close(websocket.StatusInternalError, "Internal error")
					delete(s.Clients, conn)
				}
			}
			s.Mutex.Unlock()
		}
	}
}

func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Printf("WebSocket accept error: %v", err)
		return
	}

	s.Register <- conn

	defer func() {
		s.Unregister <- conn
	}()

	for {
		var msg string
		err := wsjson.Read(context.Background(), conn, &msg)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}
	}
}

func (s *WebSocketServer) RedisSubscriber(redisClient *redis.Client) {
	var ctx = context.Background()
	pubsub := redisClient.Subscribe(ctx, "exchange_rates")
	defer pubsub.Close()

	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("Redis receive error: %v", err)
			continue
		}

		s.Broadcast <- msg.Payload
	}
}
