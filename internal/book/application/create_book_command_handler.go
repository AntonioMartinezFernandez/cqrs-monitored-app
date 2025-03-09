package book_application

import (
	"context"
	"errors"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/event"
)

type CreateBookCommandHandler struct {
	eventBus       *event.EventBus
	bookRepository book_domain.BookRepository
}

func NewCreateBookCommandHandler(
	eventBus *event.EventBus,
	bookRepository book_domain.BookRepository,
) *CreateBookCommandHandler {
	return &CreateBookCommandHandler{
		eventBus:       eventBus,
		bookRepository: bookRepository,
	}
}

func (cb CreateBookCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	cbCommand, ok := command.(*CreateBookCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	book := book_domain.NewBook(cbCommand.ID, cbCommand.Title, cbCommand.AuthorID, cbCommand.CreatedAt)
	err := cb.bookRepository.Save(ctx, *book)
	if err != nil {
		return errors.New("error saving book")
	}

	events := book.PullEvents()
	cb.eventBus.Publish(ctx, events)

	return nil
}
