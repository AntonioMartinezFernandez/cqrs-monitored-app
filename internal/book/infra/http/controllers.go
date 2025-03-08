package book_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"

	pkg_command_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/command"
	pkg_query_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/query"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
	json_api_response "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api/response"
)

func NewPostBookController(
	commandBus *pkg_command_bus.CommandBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		defer func() {
			_ = r.Body.Close()
		}()

		var createBookRequestBody CreateBookRequestBody
		err := json.Unmarshal(body, &createBookRequestBody)
		if err != nil {
			jarm.WriteErrorResponse(
				r.Context(),
				w,
				nil,
				http.StatusBadRequest,
				fmt.Errorf("invalid request body: %w", err),
			)
			return
		}

		cmd := book_application.NewCreateBookCommand(
			createBookRequestBody.ID,
			createBookRequestBody.Title,
			createBookRequestBody.Author,
			time.Now(),
		)

		err = commandBus.Dispatch(r.Context(), cmd)
		if err != nil {
			jarm.WriteErrorResponse(
				r.Context(),
				w,
				nil,
				http.StatusInternalServerError,
				fmt.Errorf("error dispatching command: %w", err),
			)
			return
		}
	}
}

func NewGetBooksController(
	queryBus *pkg_query_bus.QueryBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := book_application.NewGetAllBooksQuery()
		queryResponse, err := queryBus.Ask(r.Context(), query)

		switch err.(type) {
		case nil:
			ctx, writer := r.Context(), w
			response := queryResponse.(book_application.GetAllBooksQueryHandlerResponse)
			jarm.WriteResponse(ctx, writer, &response, http.StatusOK)
			return
		default:
			ctx, writer, errResponse := r.Context(), w, json_api_response.NewInternalServerErrorWithDetails(err.Error())
			jarm.WriteErrorResponse(ctx, writer, errResponse, http.StatusInternalServerError, err)
			return
		}
	}
}
