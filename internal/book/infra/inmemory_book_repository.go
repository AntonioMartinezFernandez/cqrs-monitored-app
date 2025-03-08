package book_infra

import (
	"context"
	"errors"
	"sync"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"
)

var _ book_domain.BookRepository = &InMemoryBookRepository{}

type InMemoryBookRepository struct {
	mut   *sync.Mutex
	books map[string]book_domain.Book
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		mut:   &sync.Mutex{},
		books: make(map[string]book_domain.Book),
	}
}

func (r *InMemoryBookRepository) FindAll(_ context.Context) ([]book_domain.Book, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	books := make([]book_domain.Book, 0, len(r.books))
	for _, book := range r.books {
		books = append(books, book)
	}
	return books, nil
}

func (r *InMemoryBookRepository) FindByID(_ context.Context, id string) (*book_domain.Book, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	book, ok := r.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return &book, nil
}

func (r *InMemoryBookRepository) Save(_ context.Context, book book_domain.Book) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.books[book.ID()] = book
	return nil
}

func (r *InMemoryBookRepository) Update(_ context.Context, book book_domain.Book) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.books[book.ID()] = book
	return nil
}

func (r *InMemoryBookRepository) Delete(_ context.Context, id string) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	delete(r.books, id)
	return nil
}
