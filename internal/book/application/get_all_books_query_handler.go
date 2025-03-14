package book_application

import (
	"context"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	pkg_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type GetAllBooksQueryHandler struct {
	bookRepository book_domain.BookRepository
}

func NewGetAllBooksQueryHandler(
	bookRepository book_domain.BookRepository,
) *GetAllBooksQueryHandler {
	return &GetAllBooksQueryHandler{
		bookRepository: bookRepository,
	}
}

func (q GetAllBooksQueryHandler) Handle(ctx context.Context, query pkg_bus.Dto) (any, error) {
	_, ok := query.(*GetAllBooksQuery)
	if !ok {
		return nil, pkg_bus.NewInvalidDto("invalid query")
	}

	books, err := q.bookRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}
