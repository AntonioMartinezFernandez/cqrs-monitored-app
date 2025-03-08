package book_application

const GetAllBooksQueryName = "GetAllBooksQuery"

type GetAllBooksQuery struct{}

func NewGetAllBooksQuery() *GetAllBooksQuery {
	return &GetAllBooksQuery{}
}

func (hq GetAllBooksQuery) Type() string {
	return GetAllBooksQueryName
}

func (hq GetAllBooksQuery) Data() map[string]any {
	return map[string]any{}
}
