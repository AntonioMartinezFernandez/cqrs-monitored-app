package book_domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBook(t *testing.T) {
	id := "1"
	title := "Test Title"
	author := "Test Author"
	createdAt := time.Now()

	book := NewBook(id, title, author, createdAt)

	assert.Equal(t, id, book.ID())
	assert.Equal(t, title, book.Title())
	assert.Equal(t, author, book.Author())
	assert.WithinDuration(t, createdAt, book.CreatedAt(), time.Second)
}

func TestBook_ID(t *testing.T) {
	id := "1"
	book := NewBook(id, "Test Title", "Test Author", time.Now())

	assert.Equal(t, id, book.ID())
}

func TestBook_Title(t *testing.T) {
	title := "Test Title"
	book := NewBook("1", title, "Test Author", time.Now())

	assert.Equal(t, title, book.Title())
}

func TestBook_Author(t *testing.T) {
	author := "Test Author"
	book := NewBook("1", "Test Title", author, time.Now())

	assert.Equal(t, author, book.Author())
}

func TestBook_CreatedAt(t *testing.T) {
	createdAt := time.Now()
	book := NewBook("1", "Test Title", "Test Author", createdAt)

	assert.WithinDuration(t, createdAt, book.CreatedAt(), time.Second)
}
