package messaging_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/messaging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ messaging.Message = &messaging.BaseMessage{}
var _ messaging.MessageRetriever = &MockMessageRetriever{}

type MockMessageRetriever struct {
	Messages chan *messaging.BaseMessage
}

func (r *MockMessageRetriever) StartRetrieving(ctx context.Context) (<-chan *messaging.BaseMessage, error) {
	return r.Messages, nil
}

func (r *MockMessageRetriever) Retrieve(ctx context.Context) (<-chan *messaging.BaseMessage, error) {
	return r.Messages, nil
}

// MockMessageHandler implements the MessageHandler interface
type MockMessageHandler struct {
	HandledMessages []messaging.Message
	Mutex           sync.Mutex
	validatorfn     messaging.ValidatorFunc
}

func (h *MockMessageHandler) HandlerFunc() (messaging.HandlerFunc, error) {
	return func(ctx context.Context, msg messaging.Message) error {
		m, ok := msg.(*MockMessage)
		if !ok {
			return fmt.Errorf("invalid message type")
		}

		if m.Key() == "" {
			return fmt.Errorf("foo is empty")
		}

		h.Mutex.Lock()
		h.HandledMessages = append(h.HandledMessages, m)
		h.Mutex.Unlock()
		return nil
	}, nil
}

func (h *MockMessageHandler) Validator() (messaging.ValidatorFunc, error) {
	if h.validatorfn != nil {
		return h.validatorfn, nil
	}

	return func(ctx context.Context, msg *messaging.BaseMessage) error {
		return nil
	}, nil
}

func (h *MockMessageHandler) MessageType() string {
	return "mock_message_type"
}

func (h *MockMessageHandler) MessageResolver() messaging.MessageResolver {
	return func(bm *messaging.BaseMessage) (messaging.Message, error) {
		return &MockMessage{
			BaseMessage: bm,
			data:        newMockData(bm.Data()),
		}, nil
	}
}

func (h *MockMessageHandler) SetValidator(validator messaging.ValidatorFunc) {
	h.validatorfn = validator
}

type MockMessageResolver struct {
}

func (m *MockMessageResolver) MessageDataFromMap(data map[string]any) (any, error) {
	if data == nil {
		return nil, fmt.Errorf("data is nil")
	}

	key, ok := data["key"].(string)
	if !ok {
		return nil, fmt.Errorf("key is not a string")
	}

	return &MockMessageData{
		Key: key,
	}, nil
}

type MockMessage struct {
	*messaging.BaseMessage

	data MockMessageData
}

func (m *MockMessage) Key() string {
	return m.data.Key
}

type MockMessageData struct {
	Key string `json:"key"`
}

func newMockData(data map[string]any) MockMessageData {
	return MockMessageData{
		Key: data["key"].(string),
	}
}

func TestMessageConsumer_Run(t *testing.T) {
	t.Parallel()
	numMessages := 100
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := logger.NewLogger("debug")
	retriever := &MockMessageRetriever{
		Messages: make(chan *messaging.BaseMessage, numMessages),
	}
	middlewares := []messaging.Middleware{
		messaging.PanicRecovererMiddleware(),
	}
	consumer := messaging.NewMessageConsumer(logger, retriever, middlewares)

	handler := &MockMessageHandler{}
	consumer.RegisterHandler(handler)

	// Start the consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run(ctx)
	}()

	// Send mock messages
	for i := range numMessages {
		msgData := map[string]any{
			"key": fmt.Sprintf("val_%d", i),
		}

		msg, err := messaging.NewBaseMessage(
			"test_id",
			"mock_message_type",
			msgData,
			"test_schema",
			map[string]any{"foo": "bar"},
			"test_subject_name",
			"test_subject_id",
			"test_external_id",
			time.Now(),
		)
		require.NoError(t, err)

		retriever.Messages <- msg
	}
	close(retriever.Messages)
	wg.Wait()

	assert.Equal(t, numMessages, len(handler.HandledMessages))
}

func TestMessageConsumer_ValidationError(t *testing.T) {
	t.Parallel()
	numMessages := 100
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger := logger.NewNullLogger()
	retriever := &MockMessageRetriever{
		Messages: make(chan *messaging.BaseMessage, numMessages),
	}
	middlewares := []messaging.Middleware{
		messaging.PanicRecovererMiddleware(),
	}
	consumer := messaging.NewMessageConsumer(logger, retriever, middlewares)

	handler := &MockMessageHandler{}
	handler.SetValidator(func(ctx context.Context, msg *messaging.BaseMessage) error {
		return fmt.Errorf("validation error")
	})
	consumer.RegisterHandler(handler)

	// Start the consumer
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.Run(ctx)
	}()

	// Send mock messages
	for i := 0; i < numMessages; i++ {
		msgData := map[string]any{
			"key": fmt.Sprintf("val_%d", i),
		}

		msg, err := messaging.NewBaseMessage(
			"test_id",
			"mock_message_type",
			msgData,
			"test_schema",
			map[string]any{"foo": "bar"},
			"test_subject_name",
			"test_subject_id",
			"test_external_id",
			time.Now(),
		)
		require.NoError(t, err)

		retriever.Messages <- msg
	}
	close(retriever.Messages)
	wg.Wait()

	assert.Equal(t, 0, len(handler.HandledMessages))
}
