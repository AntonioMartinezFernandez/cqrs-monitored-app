package book_application

import (
	"context"
	"time"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"

	pkg_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

type GetAllBooksQueryHandler struct {
	ulidProvider   pkg_utils.UlidProvider
	bookRepository book_domain.BookRepository
}

func NewGetAllBooksQueryHandler(
	ulidProvider pkg_utils.UlidProvider,
	bookRepository book_domain.BookRepository,
) *GetAllBooksQueryHandler {
	return &GetAllBooksQueryHandler{
		ulidProvider:   ulidProvider,
		bookRepository: bookRepository,
	}
}

type Book struct {
	ID        string    `json:"id" jsonapi:"primary,books"`
	Title     string    `json:"title" jsonapi:"attr,title"`
	Author    string    `json:"author" jsonapi:"attr,author"`
	CreatedAt time.Time `json:"created_at" jsonapi:"attr,created_at"`
}

type GetAllBooksQueryHandlerResponse struct {
	Id    string `json:"id" jsonapi:"primary,books"`
	Books []Book `json:"books" jsonapi:"attr,books"`
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

	var booksResponse []Book
	for _, book := range books {
		booksResponse = append(booksResponse, Book{
			ID:        book.ID(),
			Title:     book.Title(),
			Author:    book.Author(),
			CreatedAt: book.CreatedAt(),
		})
	}

	return GetAllBooksQueryHandlerResponse{
		Id:    q.ulidProvider.New().String(),
		Books: booksResponse,
	}, nil
}
