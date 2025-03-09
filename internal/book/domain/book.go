package book_domain

import "time"

type Book struct {
	id        string
	title     string
	authorID  string
	createdAt time.Time
}

func NewBook(
	id string,
	title string,
	authorID string,
	createdAt time.Time,
) *Book {
	return &Book{
		id:        id,
		title:     title,
		authorID:  authorID,
		createdAt: createdAt,
	}
}

func (b *Book) ID() string {
	return b.id
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) AuthorID() string {
	return b.authorID
}

func (b *Book) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Book) Update(
	title string,
	authorID string,
) {
	b.title = title
	b.authorID = authorID
}
