package book_application

import (
	"context"
	"errors"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type CreateBookCommandHandler struct {
	bookRepository book_domain.BookRepository
}

func NewCreateBookCommandHandler(bookRepository book_domain.BookRepository) *CreateBookCommandHandler {
	return &CreateBookCommandHandler{
		bookRepository: bookRepository,
	}
}

func (cb CreateBookCommandHandler) Handle(ctx context.Context, command bus.Dto) error {
	cbCommand, ok := command.(*CreateBookCommand)
	if !ok {
		return bus.NewInvalidDto("Invalid command")
	}

	book := book_domain.NewBook(cbCommand.ID, cbCommand.Title, cbCommand.Author, cbCommand.CreatedAt)
	err := cb.bookRepository.Save(ctx, *book)
	if err != nil {
		return errors.New("error saving book")
	}

	return nil
}
