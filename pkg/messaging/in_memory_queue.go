package messaging

import (
	"sync"

	"slices"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

type InMemoryQueue struct {
	name         string
	messageQueue []BaseMessage
	mut          *sync.Mutex
	ulidProvider utils.UlidProvider
}

func NewInMemoryQueue(name string, ulidProvider utils.UlidProvider) *InMemoryQueue {
	return &InMemoryQueue{
		name:         name,
		messageQueue: []BaseMessage{},
		mut:          &sync.Mutex{},
		ulidProvider: ulidProvider,
	}
}

func (imq *InMemoryQueue) Name() string {
	return imq.name
}

func (imq *InMemoryQueue) Publish(msg Message) error {
	imq.mut.Lock()
	defer imq.mut.Unlock()

	baseMessage, err := NewBaseMessage(
		msg.ID(),
		msg.Type(),
		msg.Data(),
		msg.DataSchema(),
		msg.Metadata(),
		msg.SubjectName(),
		msg.SubjectID(),
		imq.ulidProvider.New().String(),
		msg.CreatedAt(),
	)
	if err != nil {
		return err
	}

	imq.messageQueue = append(imq.messageQueue, *baseMessage)
	return nil
}

func (imq *InMemoryQueue) RetrieveMessages(numberOfMessages int) ([]BaseMessage, error) {
	imq.mut.Lock()
	defer imq.mut.Unlock()

	if len(imq.messageQueue) == 0 {
		return nil, nil
	}

	if numberOfMessages > len(imq.messageQueue) {
		numberOfMessages = len(imq.messageQueue)
	}

	baseMessages := imq.messageQueue[:numberOfMessages]

	return baseMessages, nil
}

func (imq *InMemoryQueue) DeleteMessage(messageID string) error {
	imq.mut.Lock()
	defer imq.mut.Unlock()

	for i, bm := range imq.messageQueue {
		if bm.ExternalID() == messageID {
			imq.messageQueue = slices.Delete(imq.messageQueue, i, i+1)
			break
		}
	}

	return nil
}

func (imq *InMemoryQueue) PendingMessages() int {
	return len(imq.messageQueue)
}
