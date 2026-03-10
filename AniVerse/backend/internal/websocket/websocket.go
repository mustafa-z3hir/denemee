package websocket

import (
	"log"
	"sync"
	"time"

	"aniverse/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Hub struct {
	clients    map[string]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	repo       *repository.Repository
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan Message
	userID string
}

type Message struct {
	Type           string `json:"type"`
	ConversationID string `json:"conversation_id"`
	SenderID       string `json:"sender_id"`
	Content        string `json:"content"`
	Timestamp      int64  `json:"timestamp"`
}

func NewHub(repo *repository.Repository) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		repo:       repo,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userID] = client
			h.mu.Unlock()
			log.Printf("User %s connected", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("User %s disconnected", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) HandleWebSocket(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	return websocket.New(func(conn *websocket.Conn) {
		client := &Client{
			hub:    h,
			conn:   conn,
			send:   make(chan Message, 256),
			userID: userID,
		}

		h.register <- client

		defer func() {
			h.unregister <- client
		}()

		go client.writePump()
		client.readPump()
	})(c)
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg.SenderID = c.userID
		msg.Timestamp = time.Now().Unix()
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteJSON(message)
		}
	}
}
