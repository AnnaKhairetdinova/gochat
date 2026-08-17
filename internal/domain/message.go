package domain

import (
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	TypeMessage MessageType = "message"
	TypeSystem  MessageType = "system"
)

type Message struct {
	UUID      uuid.UUID   `json:"uuid"`
	Room      string      `json:"room"`
	Username  string      `json:"user"`
	Text      string      `json:"text"`
	Type      MessageType `json:"type"`
	CreatedAt time.Time   `json:"timestamp"`
}

type IncomingMessage struct {
	Text string `json:"text"`
}
