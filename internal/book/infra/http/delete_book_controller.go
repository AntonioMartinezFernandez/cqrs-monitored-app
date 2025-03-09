package book_http

import (
	"fmt"
	"net/http"

	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	"github.com/gorilla/mux"

	pkg_command_bus "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/command"
	pkg_json_api "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/json-api"
)

func NewDeleteBookController(
	commandBus *pkg_command_bus.CommandBus,
	jarm *pkg_json_api.JsonApiResponseMiddleware,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bookID := mux.Vars(r)["book_id"]

		cmd := book_application.NewDeleteBookCommand(bookID)

		err := commandBus.Dispatch(r.Context(), cmd)
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
