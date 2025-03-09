package book_application

import (
	"context"
	"errors"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type DeleteBookCommandHandler struct {
	bookRepository book_domain.BookRepository
}

func NewDeleteBookCommandHandler(bookRepository book_domain.BookRepository) *DeleteBookCommandHandler {
	return &DeleteBookCommandHandler{
		bookRepository: bookRepository,
	}
}

func (cb DeleteBookCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	cbCommand, ok := command.(*DeleteBookCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	err := cb.bookRepository.Delete(ctx, cbCommand.ID)
	if err != nil {
		return errors.New("error deleting book")
	}

	return nil
}
