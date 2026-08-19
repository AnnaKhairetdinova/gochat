package ws

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/AnnaKhairetdinova/gochat/internal/broker"
	"github.com/AnnaKhairetdinova/gochat/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 4096
)

type Client struct {
	UUID     uuid.UUID
	Username string
	Room     string
	conn     *websocket.Conn
	send     chan domain.Message
	hub      *Hub
}

func NewClient(hub *Hub, conn *websocket.Conn, username, room string, bufSize int) *Client {
	return &Client{
		UUID:     uuid.New(),
		Username: username,
		Room:     room,
		conn:     conn,
		send:     make(chan domain.Message, bufSize),
		hub:      hub,
	}
}

// readPump читает сообщения от клиента в бесконечном цикле
// Запускается в отдельной горутине: go client.readPump()

func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("unexpected close error: %v", err)
			}

			break
		}

		var incoming domain.IncomingMessage

		if err := json.Unmarshal(data, &incoming); err != nil {
			log.Printf("invalid JSON from client: %v", err)
			continue
		}

		if strings.TrimSpace(incoming.Text) == "" {
			continue
		}

		msg := domain.Message{
			UUID:      uuid.New(),
			Room:      c.Room,
			Username:  c.Username,
			Text:      incoming.Text,
			Type:      domain.TypeMessage,
			CreatedAt: time.Now().UTC(),
		}

		if err := broker.Publish(ctx, msg); err != nil {
			log.Printf("failed to publish message: %v", err)
		}
	}
}

// writePump пишет сообщения клиенту и отправляет ping
// Запускается в отдельной горутине: go client.writePump()
// Это ЕДИНСТВЕННОЕ место где пишем в c.conn

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}

			data, err := json.Marshal(msg)
			if err != nil {
				log.Printf("failed to marshal message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
