package queue

import (
	"encoding/json"
	"time"
)

type Message interface {
	ID() string
	MessageType() string
	Payload() map[string]any

	MaxRetries() int
	MaxTimeout() time.Duration
}

type QueueMessage struct {
	id          string
	messageType string
	payload     map[string]any

	maxRetries int
	maxTimeout time.Duration
}

func NewQueueMessage(
	id string,
	messageType string,
	payload map[string]any,
	maxRetries int,
	maxTimeout time.Duration,
) *QueueMessage {
	return &QueueMessage{
		id:          id,
		messageType: messageType,
		payload:     payload,

		maxRetries: maxRetries,
		maxTimeout: maxTimeout,
	}
}

func (m *QueueMessage) ID() string {
	return m.id
}

func (m *QueueMessage) MessageType() string {
	return m.messageType
}

func (m *QueueMessage) Payload() map[string]any {
	return m.payload
}

func (m *QueueMessage) MaxRetries() int {
	return m.maxRetries
}

func (m *QueueMessage) MaxTimeout() time.Duration {
	return m.maxTimeout
}

func (m *QueueMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          string         `json:"id"`
		MessageType string         `json:"messageType"`
		Payload     map[string]any `json:"payload"`
		MaxRetries  int            `json:"maxRetries"`
		MaxTimeout  time.Duration  `json:"maxTimeout"`
	}{
		ID:          m.id,
		MessageType: m.messageType,
		Payload:     m.payload,
		MaxRetries:  m.maxRetries,
		MaxTimeout:  m.maxTimeout,
	})
}

func (m *QueueMessage) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID          string         `json:"id"`
		MessageType string         `json:"messageType"`
		Payload     map[string]any `json:"payload"`
		MaxRetries  int            `json:"maxRetries"`
		MaxTimeout  time.Duration  `json:"maxTimeout"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	m.id = aux.ID
	m.messageType = aux.MessageType
	m.payload = aux.Payload
	m.maxRetries = aux.MaxRetries
	m.maxTimeout = aux.MaxTimeout

	return nil
}
