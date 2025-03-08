package book_http

import (
	"net/http"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"
	pkg_query_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/query"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api/response"
	pkg_utils "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/utils"
	"github.com/gorilla/mux"
)

func NewGetBookController(
	ulidProvider pkg_utils.UlidProvider,
	queryBus *pkg_query_bus.QueryBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := mux.Vars(r)["book_id"]

		query := book_application.NewGetBookByIDQuery(bookID)
		queryResponse, err := queryBus.Ask(r.Context(), query)

		switch err.(type) {
		case nil:
			ctx, writer := r.Context(), w
			book, _ := queryResponse.(*book_domain.Book)

			bookResponse := Book{
				ID:        book.ID(),
				Title:     book.Title(),
				Author:    book.Author(),
				CreatedAt: book.CreatedAt(),
			}

			response := GetBookByIDQueryHandlerResponse{
				Id:   ulidProvider.New().String(),
				Book: bookResponse,
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
