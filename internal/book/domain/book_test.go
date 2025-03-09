package book_domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBook(t *testing.T) {
	id := "01JNX82KTP44N21RXD4N0G9RBG"
	title := "Test Title"
	authorID := "01JNX81X6JYYC3E2P3B89Q7WW5"
	createdAt := time.Now()

	book := NewBook(id, title, authorID, createdAt)

	assert.Equal(t, id, book.ID())
	assert.Equal(t, title, book.Title())
	assert.Equal(t, authorID, book.AuthorID())
	assert.WithinDuration(t, createdAt, book.CreatedAt(), time.Second)
}

func TestBook_ID(t *testing.T) {
	id := "01JNX82KTP44N21RXD4N0G9RBG"
	book := NewBook(id, "Test Title", "01JNX81X6JYYC3E2P3B89Q7WW5", time.Now())

	assert.Equal(t, id, book.ID())
}

func TestBook_Title(t *testing.T) {
	title := "Test Title"
	book := NewBook("01JNX82KTP44N21RXD4N0G9RBG", title, "01JNX81X6JYYC3E2P3B89Q7WW5", time.Now())

	assert.Equal(t, title, book.Title())
}

func TestBook_AuthorID(t *testing.T) {
	authorID := "01JNX81X6JYYC3E2P3B89Q7WW5"
	book := NewBook("01JNX82KTP44N21RXD4N0G9RBG", "Test Title", authorID, time.Now())

	assert.Equal(t, authorID, book.AuthorID())
}

func TestBook_CreatedAt(t *testing.T) {
	createdAt := time.Now()
	book := NewBook("01JNX82KTP44N21RXD4N0G9RBG", "Test Title", "01JNX81X6JYYC3E2P3B89Q7WW5", createdAt)

	assert.WithinDuration(t, createdAt, book.CreatedAt(), time.Second)
}

func TestBook_Update(t *testing.T) {
	book := NewBook("01JNX82KTP44N21RXD4N0G9RBG", "Test Title", "01JNX81X6JYYC3E2P3B89Q7WW5", time.Now())

	newTitle := "New Title"
	newAuthorID := "this-is-a-new-author-id"
	book.Update(newTitle, newAuthorID)

	assert.Equal(t, newTitle, book.Title())
	assert.Equal(t, newAuthorID, book.AuthorID())
}
