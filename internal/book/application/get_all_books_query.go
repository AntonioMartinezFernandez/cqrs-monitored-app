package book_application

const getAllBooksQueryName = "GetAllBooksQuery"

type GetAllBooksQuery struct{}

func NewGetAllBooksQuery() *GetAllBooksQuery {
	return &GetAllBooksQuery{}
}

func (hq GetAllBooksQuery) Type() string {
	return getAllBooksQueryName
}

func (hq GetAllBooksQuery) Data() map[string]any {
	return map[string]any{}
}
