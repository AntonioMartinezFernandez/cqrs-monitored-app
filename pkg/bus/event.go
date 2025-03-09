package bus

import "encoding/json"

type Event interface {
	Dto

	ID() string
	Name() string
	Schema() string
	Data() map[string]any
	MetaData() map[string]any

	Serialize() []byte
	Deserialize([]byte) Event
}

type Events []Event

type BaseEvent struct {
	EventId       string         `json:"id"`
	EventName     string         `json:"name"`
	EventType     string         `json:"type"`
	EventSchema   string         `json:"schema"`
	EventData     map[string]any `json:"data"`
	EventMetaData map[string]any `json:"metadata"`
}

func (e *BaseEvent) Type() string {
	return e.EventType
}

func (e *BaseEvent) ID() string {
	return e.EventId
}

func (e *BaseEvent) Name() string {
	return e.EventName
}

func (e *BaseEvent) Schema() string {
	return e.EventSchema
}

func (e *BaseEvent) Data() map[string]any {
	return e.EventData
}

func (e *BaseEvent) MetaData() map[string]any {
	return e.EventMetaData
}

func (e *BaseEvent) Serialize() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return data
}

func (e *BaseEvent) Deserialize(data []byte) Event {
	var event BaseEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil
	}
	return &event
}
