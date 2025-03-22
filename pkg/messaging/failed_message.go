package messaging

import (
	"time"
)

const (
	PendingStatus   FailedMessageStatus = "pending"
	PublishedStatus FailedMessageStatus = "published"
)

// FailedMessageStatus defines whether a message has been published or not
type FailedMessageStatus string

func (fms FailedMessageStatus) String() string {
	return string(fms)
}

// FailedMessage represents an unpublished message given a failure or so
type FailedMessage struct {
	id       string
	payload  []byte
	status   FailedMessageStatus
	failedAt time.Time
}

// NewFailedMessage creates a new FailedMessage as pending
func NewFailedMessage(id string, payload []byte, failedAt time.Time) *FailedMessage {
	return &FailedMessage{
		id:       id,
		payload:  payload,
		status:   PendingStatus,
		failedAt: failedAt,
	}
}

// NewFailedMessageFromPrimitives creates a new FailedMessage from
// its primitives values
func NewFailedMessageFromPrimitives(
	id string,
	payload []byte,
	status string,
	failedAt time.Time,
) *FailedMessage {
	return &FailedMessage{
		id:       id,
		payload:  payload,
		status:   FailedMessageStatus(status),
		failedAt: failedAt,
	}
}

func (fm *FailedMessage) Identifier() string {
	return fm.id
}

func (fm *FailedMessage) Payload() []byte {
	return fm.payload
}

func (fm *FailedMessage) Status() FailedMessageStatus {
	return fm.status
}

func (fm *FailedMessage) FailedAt() time.Time {
	return fm.failedAt
}

// MarkAsPublished marks the message as published on its destination
func (fm *FailedMessage) MarkAsPublished() {
	fm.status = PublishedStatus
}
