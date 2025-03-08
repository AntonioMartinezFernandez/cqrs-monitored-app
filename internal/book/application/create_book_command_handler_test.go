package book_application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	book_domain_mocks "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CreateBookCommand is a mock implementation of the bus.Dto interface
type CreateBookCommand struct {
	ID        string
	Title     string
	Author    string
	CreatedAt time.Time
}

func TestCreateBookCommandHandler_Handle(t *testing.T) {

	t.Run("successfully create book", func(t *testing.T) {
		ctx := context.Background()

		mockRepo := book_domain_mocks.NewBookRepository(t)
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil).Once()

		handler := book_application.NewCreateBookCommandHandler(mockRepo)

		cmd := book_application.NewCreateBookCommand("1", "Test Book", "Test Author", time.Now())

		err := handler.Handle(ctx, cmd)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("error saving book", func(t *testing.T) {
		ctx := context.Background()

		mockRepo := book_domain_mocks.NewBookRepository(t)
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("error saving book")).Once()

		handler := book_application.NewCreateBookCommandHandler(mockRepo)

		cmd := book_application.NewCreateBookCommand("1", "Test Book", "Test Author", time.Now())

		err := handler.Handle(ctx, cmd)
		assert.Error(t, err)
		assert.Equal(t, "error saving book", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
