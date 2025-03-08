package book_http

type CreateBookRequestBody struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
