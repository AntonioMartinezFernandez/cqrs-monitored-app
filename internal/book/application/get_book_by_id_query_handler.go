package book_application

import (
	"context"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	pkg_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type GetBookByIDQueryHandler struct {
	bookRepository book_domain.BookRepository
}

func NewGetBookByIDQueryHandler(
	bookRepository book_domain.BookRepository,
) *GetBookByIDQueryHandler {
	return &GetBookByIDQueryHandler{
		bookRepository: bookRepository,
	}
}

func (q GetBookByIDQueryHandler) Handle(ctx context.Context, query pkg_bus.Dto) (any, error) {
	parsedQuery, ok := query.(*GetBookByIDQuery)
	if !ok {
		return nil, pkg_bus.NewInvalidDto("invalid query")
	}

	book, err := q.bookRepository.FindByID(ctx, parsedQuery.ID)
	if err != nil {
		return nil, err
	}

	return book, nil
}
