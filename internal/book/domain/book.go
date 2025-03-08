package book_domain

import "time"

type Book struct {
	id        string
	title     string
	author    string
	createdAt time.Time
}

func NewBook(
	id string,
	title string,
	author string,
	createdAt time.Time,
) *Book {
	return &Book{
		id:        id,
		title:     title,
		author:    author,
		createdAt: createdAt,
	}
}

func (b *Book) ID() string {
	return b.id
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) Author() string {
	return b.author
}

func (b *Book) CreatedAt() time.Time {
	return b.createdAt
}
