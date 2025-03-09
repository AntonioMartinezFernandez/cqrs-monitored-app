package author_domain

import "context"

type AuthorRepository interface {
	FindAll(ctx context.Context) ([]Author, error)
	FindByID(ctx context.Context, id string) (*Author, error)
	Save(ctx context.Context, author Author) error
	Update(ctx context.Context, author Author) error
	Delete(ctx context.Context, id string) error
}
