package book_domain

import "context"

type BookRepository interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindByID(ctx context.Context, id string) (*Book, error)
	Save(ctx context.Context, book Book) error
	Update(ctx context.Context, book Book) error
	Delete(ctx context.Context, id string) error
}
