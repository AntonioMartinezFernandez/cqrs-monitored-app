package messaging_test

import (
	"testing"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/messaging"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryQueue(t *testing.T) {
	msg1, err := NewExampleMessage("message-id-1", "Michael Jordan", 95, map[string]any{}, time.Now())
	assert.NoError(t, err)

	msg2, err := NewExampleMessage("message-id-2", "KobeBryant", 94, map[string]any{}, time.Now())
	assert.NoError(t, err)

	queue := messaging.NewInMemoryQueue("test-queue", utils.NewRandomUlidProvider())
	queue.Publish(msg1)
	queue.Publish(msg2)

	availableMessages := queue.PendingMessages()
	assert.Equal(t, availableMessages, 2)

	retrievedMsg, err := queue.RetrieveMessages(1)
	assert.NoError(t, err)
	assert.Equal(t, len(retrievedMsg), 1)
	assert.Equal(t, retrievedMsg[0].ID(), "message-id-1")

	queue.DeleteMessage(retrievedMsg[0].ExternalID())
	availableMessages = queue.PendingMessages()
	assert.Equal(t, availableMessages, 1)
}

const ExampleMessageType = "company.example-entity.example-message"
const ExampleMessageSchema = "schemas/example-message.schema.json"

type ExampleMessage struct {
	*messaging.BaseMessage

	playerName  string
	playerScore int
}

func (m *ExampleMessage) Name() string {
	return m.playerName
}

func (m *ExampleMessage) Age() int {
	return m.playerScore
}

func NewExampleMessage(
	id string,
	playerName string,
	playerScore int,
	msgMetadata map[string]any,
	msgCreatedAt time.Time,
) (*ExampleMessage, error) {
	msgData := map[string]any{
		"playerName":  playerName,
		"playerScore": playerScore,
	}

	baseMessage, err := messaging.NewBaseMessage(
		id,
		ExampleMessageType,
		msgData,
		ExampleMessageSchema,
		msgMetadata,
		"example-entity",
		"example-entity-id",
		id,
		msgCreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ExampleMessage{
		BaseMessage: baseMessage,

		playerName:  playerName,
		playerScore: playerScore,
	}, nil
}
