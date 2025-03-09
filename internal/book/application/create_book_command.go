package book_application

import "time"

const createBookCommandName = "create_book_command"
const createBookBlockingKey = "book_operation"

type CreateBookCommand struct {
	ID        string
	Title     string
	AuthorID  string
	CreatedAt time.Time
}

func NewCreateBookCommand(
	id string,
	title string,
	authorID string,
	createdAt time.Time,
) *CreateBookCommand {
	return &CreateBookCommand{
		ID:        id,
		Title:     title,
		AuthorID:  authorID,
		CreatedAt: createdAt,
	}
}

func (cdp *CreateBookCommand) Type() string {
	return createBookCommandName
}

func (cdp *CreateBookCommand) BlockingKey() string {
	return createBookBlockingKey
}
