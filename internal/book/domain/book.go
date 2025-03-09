package book_domain

import (
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/domain"
)

type Book struct {
	id        string
	title     string
	authorID  string
	createdAt time.Time

	recorder *domain.EventRecorder
}

func NewBook(
	id string,
	title string,
	authorID string,
	createdAt time.Time,
) *Book {
	bookCreatedEvent := NewBookCreatedEvent(
		id,
		title,
		authorID,
		createdAt,
	)

	book := &Book{
		id:        id,
		title:     title,
		authorID:  authorID,
		createdAt: createdAt,

		recorder: domain.NewEventRecorder(),
	}

	book.recorder.Record(bookCreatedEvent)

	return book
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

func (b *Book) PullEvents() bus.Events {
	return b.recorder.Pull()
}

func (b *Book) Update(
	title string,
	authorID string,
) {
	b.title = title
	b.authorID = authorID
}
