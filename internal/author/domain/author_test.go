package author_domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthor(t *testing.T) {
	author := NewAuthor("1", "John Doe")

	assert.Equal(t, "1", author.ID())
	assert.Equal(t, "John Doe", author.Name())
}

func TestAuthor_Update(t *testing.T) {
	author := NewAuthor("1", "John Doe")
	author.Update("Jane Doe")

	assert.Equal(t, "Jane Doe", author.Name())
}
