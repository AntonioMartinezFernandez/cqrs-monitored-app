package book_application

const updateBookCommandName = "update_book_command"
const updateBookBlockingKey = "book_operation"

type UpdateBookCommand struct {
	ID       string
	Title    string
	AuthorID string
}

func NewUpdateBookCommand(
	id string,
	title string,
	authorID string,
) *UpdateBookCommand {
	return &UpdateBookCommand{
		ID:       id,
		Title:    title,
		AuthorID: authorID,
	}
}

func (cdp *UpdateBookCommand) Type() string {
	return updateBookCommandName
}

func (cdp *UpdateBookCommand) BlockingKey() string {
	return updateBookBlockingKey
}
