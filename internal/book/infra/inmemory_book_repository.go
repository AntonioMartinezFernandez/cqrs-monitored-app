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
	books []book_domain.Book
}

func NewInMemoryBookRepository() *InMemoryBookRepository {
	return &InMemoryBookRepository{
		mut:   &sync.Mutex{},
		books: make([]book_domain.Book, 0),
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

	for _, book := range r.books {
		if book.ID() == id {
			return &book, nil
		}
	}

	return nil, errors.New("book not found")
}

func (r *InMemoryBookRepository) Save(_ context.Context, newBook book_domain.Book) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	for _, book := range r.books {
		if book.ID() == newBook.ID() {
			return errors.New("book already exists")
		}
	}

	r.books = append(r.books, newBook)

	return nil
}

func (r *InMemoryBookRepository) Update(_ context.Context, book book_domain.Book) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	for i, oldBook := range r.books {
		if oldBook.ID() == book.ID() {
			oldBook.Update(book.Title(), book.AuthorID())
			r.books[i] = oldBook
			return nil
		}
	}

	return errors.New("book not found")
}

func (r *InMemoryBookRepository) Delete(_ context.Context, id string) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	books := make([]book_domain.Book, 0, len(r.books))

	for _, book := range r.books {
		if book.ID() != id {
			books = append(books, book)
		}
	}

	r.books = books

	return nil
}
