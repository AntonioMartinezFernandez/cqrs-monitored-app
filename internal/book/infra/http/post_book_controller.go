package book_http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"

	pkg_command_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/command"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
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
