package book_application

const deleteBookCommandName = "delete_book_command"
const deleteBookBlockingKey = "book_operation"

type DeleteBookCommand struct {
	ID string
}

func NewDeleteBookCommand(
	id string,
) *DeleteBookCommand {
	return &DeleteBookCommand{
		ID: id,
	}
}

func (cdp *DeleteBookCommand) Type() string {
	return deleteBookCommandName
}

func (cdp *DeleteBookCommand) BlockingKey() string {
	return deleteBookBlockingKey
}
