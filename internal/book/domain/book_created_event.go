package book_domain

import (
	"encoding/json"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

const (
	BookCreatedEventName   = "bookCreatedV1"
	bookCreatedEventType   = "book.created.v1.event"
	bookCreatedEventSchema = "book-created.schema.json"
)

type BookCreatedEvent struct {
	EventId       string         `json:"id"`
	EventName     string         `json:"name"`
	EventType     string         `json:"type"`
	EventSchema   string         `json:"schema"`
	EventData     map[string]any `json:"data"`
	EventMetaData map[string]any `json:"metadata"`
}

func NewBookCreatedEvent(
	bookID string,
	bookName string,
	authorID string,
	createdAt time.Time,
) *BookCreatedEvent {
	data := map[string]any{
		"id":        bookID,
		"name":      bookName,
		"authorID":  authorID,
		"createdAt": createdAt,
	}

	return &BookCreatedEvent{
		EventId:       pkg_utils.NewUlid().String(),
		EventName:     BookCreatedEventName,
		EventType:     bookCreatedEventType,
		EventSchema:   bookCreatedEventSchema,
		EventData:     data,
		EventMetaData: map[string]any{},
	}
}

func (bc *BookCreatedEvent) Type() string {
	return bc.EventType
}
func (bc *BookCreatedEvent) ID() string {
	return bc.EventId
}
func (bc *BookCreatedEvent) Name() string {
	return bc.EventName
}
func (bc *BookCreatedEvent) Schema() string {
	return bc.EventSchema
}
func (bc *BookCreatedEvent) Data() map[string]any {
	return bc.EventData
}
func (bc *BookCreatedEvent) MetaData() map[string]any {
	return bc.EventMetaData
}
func (bc *BookCreatedEvent) Serialize() []byte {
	data, err := json.Marshal(bc)
	if err != nil {
		return nil
	}
	return data
}
func (bc *BookCreatedEvent) Deserialize(data []byte) bus.Event {
	var event BookCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil
	}
	return &event
}
