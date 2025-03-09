package book_application

import (
	"context"
	"errors"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type UpdateBookCommandHandler struct {
	bookRepository book_domain.BookRepository
}

func NewUpdateBookCommandHandler(bookRepository book_domain.BookRepository) *UpdateBookCommandHandler {
	return &UpdateBookCommandHandler{
		bookRepository: bookRepository,
	}
}

func (cb UpdateBookCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	cbCommand, ok := command.(*UpdateBookCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	book, err := cb.bookRepository.FindByID(ctx, cbCommand.ID)
	if err != nil {
		return errors.New("error finding book")
	}

	book.Update(
		cbCommand.Title,
		cbCommand.AuthorID,
	)

	updErr := cb.bookRepository.Update(ctx, *book)
	if updErr != nil {
		return errors.New("error saving book")
	}

	return nil
}
