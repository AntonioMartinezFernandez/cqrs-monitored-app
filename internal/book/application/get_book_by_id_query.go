package book_application

const getBookByIDQueryName = "GetBookByIDQuery"

type GetBookByIDQuery struct {
	ID string
}

func NewGetBookByIDQuery(id string) *GetBookByIDQuery {
	return &GetBookByIDQuery{
		ID: id,
	}
}

func (hq GetBookByIDQuery) Type() string {
	return getBookByIDQueryName
}

func (hq GetBookByIDQuery) Data() map[string]any {
	return map[string]any{}
}
