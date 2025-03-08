package book_application

const GetBookByIDQueryName = "GetBookByIDQuery"

type GetBookByIDQuery struct {
	ID string
}

func NewGetBookByIDQuery(id string) *GetBookByIDQuery {
	return &GetBookByIDQuery{
		ID: id,
	}
}

func (hq GetBookByIDQuery) Type() string {
	return GetBookByIDQueryName
}

func (hq GetBookByIDQuery) Data() map[string]any {
	return map[string]any{}
}
