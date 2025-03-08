package book_http

import (
	"net/http"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"
	pkg_query_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/query"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api/response"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
)

func NewGetBooksController(
	ulidProvider pkg_utils.UlidProvider,
	queryBus *pkg_query_bus.QueryBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := book_application.NewGetAllBooksQuery()
		queryResponse, err := queryBus.Ask(r.Context(), query)

		switch err.(type) {
		case nil:
			ctx, writer := r.Context(), w
			books, _ := queryResponse.([]book_domain.Book)

			var booksResponse []Book
			for _, book := range books {
				booksResponse = append(booksResponse, Book{
					ID:        book.ID(),
					Title:     book.Title(),
					Author:    book.Author(),
					CreatedAt: book.CreatedAt(),
				})
			}

			response := GetAllBooksQueryHandlerResponse{
				Id:    ulidProvider.New().String(),
				Books: booksResponse,
			}

			jarm.WriteResponse(ctx, writer, &response, http.StatusOK)
			return
		default:
			ctx, writer, errResponse := r.Context(), w, json_api_response.NewInternalServerErrorWithDetails(err.Error())
			jarm.WriteErrorResponse(ctx, writer, errResponse, http.StatusInternalServerError, err)
			return
		}
	}
}
