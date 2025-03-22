package messaging_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/messaging"
	"github.com/stretchr/testify/assert"
)

func TestBaseMessage_Nack(t *testing.T) {
	// Create a new BaseMessage instance
	msg, err := messaging.NewBaseMessage(
		"test_id",
		"test_type",
		map[string]any{"key": "value"},
		"test_schema",
		map[string]any{"foo": "bar"},
		"test_subject_name",
		"test_subject_id",
		"test_external_id",
		time.Now(),
	)
	assert.NoError(t, err)

	go func() {
		// Call Nack to mark the message as not acknowledged
		msg.Nack()
	}()

	for {
		select {
		case <-msg.Nacked():
			return
		case <-time.After(5 * time.Second):
			t.Fatal("Nacked channel should be closed after calling Nack")
			return
		}
	}
}

func TestBaseMessage_Ack(t *testing.T) {
	// Create a new BaseMessage instance
	msg, err := messaging.NewBaseMessage(
		"test_id",
		"test_type",
		map[string]any{"key": "value"},
		"test_schema",
		map[string]any{"foo": "bar"},
		"test_subject_name",
		"test_subject_id",
		"test_external_id",
		time.Now(),
	)
	assert.NoError(t, err)

	go func() {
		// Call Ack to mark the message as acknowledged
		msg.Ack()
	}()

	for {
		select {
		case <-msg.Acked():
			return
		case <-time.After(5 * time.Second):
			t.Fatal("Acked channel should be closed after calling Ack")
			return
		}
	}
}

func TestBaseMessage_SetContext(t *testing.T) {
	// Create a new BaseMessage instance
	msg, err := messaging.NewBaseMessage(
		"test_id",
		"test_type",
		map[string]any{"key": "value"},
		"test_schema",
		map[string]any{"foo": "bar"},
		"test_subject_name",
		"test_subject_id",
		"test_external_id",
		time.Now(),
	)
	assert.NoError(t, err)

	// Create a new context
	type contextKey string // Define a custom type for the context key
	ctx := context.WithValue(context.Background(), contextKey("key"), "value")

	// Set the context of the message
	msg.SetContext(ctx)

	// Verify that the context is set correctly
	assert.Equal(t, ctx, msg.Context())
}

func TestBaseMessage_MarshalUnmarshalJSON(t *testing.T) {
	// Create a new BaseMessage instance
	originalMsg, err := messaging.NewBaseMessage(
		"test_id",
		"test_type",
		map[string]any{"key": "value"},
		"test_schema",
		map[string]any{"foo": "bar"},
		"test_subject_name",
		"test_subject_id",
		"test_external_id",
		time.Now(),
	)
	assert.NoError(t, err)

	// Marshal the message to JSON
	data, err := json.Marshal(originalMsg)
	assert.NoError(t, err)

	// Unmarshal the JSON back to BaseMessage
	var unmarshaledMsg messaging.BaseMessage
	err = json.Unmarshal(data, &unmarshaledMsg)
	assert.NoError(t, err)

	// Compare the original and unmarshaled messages
	assert.Equal(t, originalMsg.ID(), unmarshaledMsg.ID())
	assert.Equal(t, originalMsg.Type(), unmarshaledMsg.Type())
	assert.Equal(t, originalMsg.Data(), unmarshaledMsg.Data())
	assert.Equal(t, originalMsg.DataSchema(), unmarshaledMsg.DataSchema())
	assert.Equal(t, originalMsg.Metadata(), unmarshaledMsg.Metadata())
	assert.Equal(t, originalMsg.SubjectName(), unmarshaledMsg.SubjectName())
	assert.Equal(t, originalMsg.SubjectID(), unmarshaledMsg.SubjectID())
	assert.Equal(t, originalMsg.ExternalID(), unmarshaledMsg.ExternalID())
	assert.Equal(t, originalMsg.CreatedAt(), unmarshaledMsg.CreatedAt())
}

func TestBaseMessage_Getters(t *testing.T) {
	// Create a new BaseMessage instance
	msgTime := time.Now().UTC()
	msg, err := messaging.NewBaseMessage(
		"test_id",
		"test_type",
		map[string]any{"key": "value"},
		"test_schema",
		map[string]any{"foo": "bar"},
		"test_subject_name",
		"test_subject_id",
		"test_external_id",
		msgTime,
	)
	assert.NoError(t, err)

	// Test getter methods
	assert.Equal(t, "test_id", msg.ID())
	assert.Equal(t, "test_type", msg.Type())
	assert.Equal(t, map[string]any{"key": "value"}, msg.Data())
	assert.Equal(t, "test_schema", msg.DataSchema())
	assert.Equal(t, map[string]any{"foo": "bar"}, msg.Metadata())
	assert.Equal(t, "test_subject_name.test_subject_id", msg.Subject())
	assert.Equal(t, "test_subject_name", msg.SubjectName())
	assert.Equal(t, "test_subject_id", msg.SubjectID())
	assert.Equal(t, "test_external_id", msg.ExternalID())
	assert.Equal(t, msgTime, msg.CreatedAt())
}
