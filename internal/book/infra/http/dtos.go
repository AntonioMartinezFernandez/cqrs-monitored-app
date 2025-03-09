package book_http

import "time"

type CreateBookRequestBody struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	AuthorID string `json:"authorID"`
}

type Book struct {
	ID        string    `json:"id" jsonapi:"primary,books"`
	Title     string    `json:"title" jsonapi:"attr,title"`
	AuthorID  string    `json:"authorID" jsonapi:"attr,authorID"`
	CreatedAt time.Time `json:"created_at" jsonapi:"attr,created_at"`
}

type GetBookByIDQueryHandlerResponse struct {
	Id   string `json:"id" jsonapi:"primary,books"`
	Book Book   `json:"book" jsonapi:"attr,book"`
}

type GetAllBooksQueryHandlerResponse struct {
	Id    string `json:"id" jsonapi:"primary,books"`
	Books []Book `json:"books" jsonapi:"attr,books"`
}
