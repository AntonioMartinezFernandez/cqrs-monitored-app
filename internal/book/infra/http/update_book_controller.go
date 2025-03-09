package book_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	"github.com/gorilla/mux"

	pkg_command_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/command"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
)

func NewPutBookController(
	commandBus *pkg_command_bus.CommandBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := mux.Vars(r)["book_id"]

		body, _ := io.ReadAll(r.Body)
		defer func() {
			_ = r.Body.Close()
		}()

		var updateBookRequestBody UpdateBookRequestBody
		err := json.Unmarshal(body, &updateBookRequestBody)
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

		if bookID != updateBookRequestBody.ID {
			jarm.WriteErrorResponse(
				r.Context(),
				w,
				nil,
				http.StatusBadRequest,
				fmt.Errorf("invalid request body: %w", err),
			)
			return
		}

		cmd := book_application.NewUpdateBookCommand(
			updateBookRequestBody.ID,
			updateBookRequestBody.Title,
			updateBookRequestBody.AuthorID,
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
