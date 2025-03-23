package messaging

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

var _ bus.Dto = (*BaseMessage)(nil)

type BaseMessage struct {
	msgId          string
	msgType        string
	msgData        map[string]any
	msgDataSchema  string
	msgMetadata    map[string]any
	msgSubjectName string
	msgSubjectID   string
	msgExternalID  string
	msgCreatedAt   time.Time

	ctx      context.Context
	ack      chan struct{}
	noAck    chan struct{}
	ackOnce  *sync.Once
	nackOnce *sync.Once
}

func NewBaseMessage(
	msgId string,
	msgType string,
	msgData map[string]any,
	msgDataSchema string,
	msgMetadata map[string]any,
	msgSubjectName string,
	msgSubjectID string,
	msgExternalID string,
	msgCreatedAt time.Time,
) (*BaseMessage, error) {
	return &BaseMessage{
		msgId:          msgId,
		msgType:        msgType,
		msgData:        msgData,
		msgDataSchema:  msgDataSchema,
		msgMetadata:    msgMetadata,
		msgSubjectName: msgSubjectName,
		msgSubjectID:   msgSubjectID,
		msgExternalID:  msgExternalID,
		msgCreatedAt:   msgCreatedAt.UTC(),

		ack:      make(chan struct{}),
		noAck:    make(chan struct{}),
		ackOnce:  &sync.Once{},
		nackOnce: &sync.Once{},
	}, nil
}

func (bm *BaseMessage) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ID          string         `json:"id"`
		MessageType string         `json:"messageType"`
		Data        map[string]any `json:"data"`
		DataSchema  string         `json:"dataSchema"`
		Metadata    map[string]any `json:"metadata"`
		SubjectName string         `json:"subjectName"`
		SubjectID   string         `json:"subjectID"`
		ExternalID  string         `json:"externalID"`
		CreatedAt   time.Time      `json:"createdAt"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	bm.msgId = aux.ID
	bm.msgType = aux.MessageType
	bm.msgData = aux.Data
	bm.msgDataSchema = aux.DataSchema
	bm.msgMetadata = aux.Metadata
	bm.msgSubjectName = aux.SubjectName
	bm.msgSubjectID = aux.SubjectID
	bm.msgExternalID = aux.ExternalID
	bm.msgCreatedAt = aux.CreatedAt

	bm.ack = make(chan struct{})
	bm.noAck = make(chan struct{})
	bm.ackOnce = &sync.Once{}
	bm.nackOnce = &sync.Once{}

	return nil
}

func (bm *BaseMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID          string         `json:"id"`
		MessageType string         `json:"messageType"`
		Data        map[string]any `json:"data"`
		DataSchema  string         `json:"dataSchema"`
		Metadata    map[string]any `json:"metadata"`
		SubjectName string         `json:"subjectName"`
		SubjectID   string         `json:"subjectID"`
		ExternalID  string         `json:"externalID"`
		CreatedAt   time.Time      `json:"createdAt"`
	}{
		ID:          bm.msgId,
		MessageType: bm.msgType,
		Data:        bm.msgData,
		DataSchema:  bm.msgDataSchema,
		Metadata:    bm.msgMetadata,
		SubjectName: bm.msgSubjectName,
		SubjectID:   bm.msgSubjectID,
		ExternalID:  bm.msgExternalID,
		CreatedAt:   bm.msgCreatedAt,
	})
}

// ID returns the unique identifier of the message
func (bm *BaseMessage) ID() string {
	return bm.msgId
}

// Type returns the type of the message
func (bm *BaseMessage) Type() string {
	return bm.msgType
}

// Data returns the payload data of the message as map
func (bm *BaseMessage) Data() map[string]any {
	return bm.msgData
}

// DataSchema returns the JSON schema relative path of the message
func (bm *BaseMessage) DataSchema() string {
	return bm.msgDataSchema
}

// Metadata returns the metadata of the message
func (bm *BaseMessage) Metadata() map[string]any {
	return bm.msgMetadata
}

// SubjectName returns the subject name of the message
func (bm *BaseMessage) SubjectName() string {
	return bm.msgSubjectName
}

// SubjectID returns the subject id of the message
func (bm *BaseMessage) SubjectID() string {
	return bm.msgSubjectID
}

// Subject returns the subject of the message (with format [subject_name].[subject_id])
func (bm *BaseMessage) Subject() string {
	return bm.msgSubjectName + "." + bm.msgSubjectID
}

// ExternalID returns the identifier of the external message queue system
func (bm *BaseMessage) ExternalID() string {
	return bm.msgExternalID
}

// Time returns the time when the message was created
func (bm *BaseMessage) CreatedAt() time.Time {
	return bm.msgCreatedAt
}

// Ack marks the message as acknowledged
func (bm *BaseMessage) Ack() {
	bm.ackOnce.Do(func() {
		select {
		case <-bm.ack:
			// Channel is already closed
		default:
			close(bm.ack)
		}
	})
}

// Nack marks the message as not acknowledged
func (bm *BaseMessage) Nack() {
	bm.nackOnce.Do(func() {
		select {
		case <-bm.noAck:
			// Channel is already closed
		default:
			close(bm.noAck)
		}
	})
}

// Acked returns a channel that is closed when the message is acknowledged
func (bm *BaseMessage) Acked() <-chan struct{} {
	return bm.ack
}

// Nacked returns a channel that is closed when the message is not acknowledged
func (bm *BaseMessage) Nacked() <-chan struct{} {
	return bm.noAck
}

// Context returns the context of the message.
// If the context is not set, it returns a background context.
// This is useful for the handlers that wants to introduce data in the context.
func (bm *BaseMessage) Context() context.Context {
	if bm.ctx != nil {
		return bm.ctx
	}
	return context.Background()
}

// SetContext sets the context of the message
func (bm *BaseMessage) SetContext(ctx context.Context) {
	bm.ctx = ctx
}
