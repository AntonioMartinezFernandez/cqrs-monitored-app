package messaging

import (
	"time"
)

type Message interface {
	ID() string
	Type() string
	Data() map[string]any
	DataSchema() string
	Metadata() map[string]any
	SubjectName() string
	SubjectID() string
	CreatedAt() time.Time
}
