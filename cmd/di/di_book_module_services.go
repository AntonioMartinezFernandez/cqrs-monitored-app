package di

import (
	book_application "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/application"
	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"
	book_infra "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/infra"
	book_http "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/infra/http"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus/event"
)

type BookModuleServices struct {
	CreateBookCommandHandler *book_application.CreateBookCommandHandler
	DeleteBookCommandHandler *book_application.DeleteBookCommandHandler
	UpdateBookCommandHandler *book_application.UpdateBookCommandHandler
	GetAllBooksQueryHandler  *book_application.GetAllBooksQueryHandler
	GetBookByIDQueryHandler  *book_application.GetBookByIDQueryHandler

	PrintBookOnBookCreatedEventHandler event.EventHandler
}

func InitBookModuleServices(commonServices *CommonServices, httpServices *HttpServices) *BookModuleServices {
	bookRepository := book_infra.NewInMemoryBookRepository()

	createBookCommandHandler := book_application.NewCreateBookCommandHandler(commonServices.EventBus, bookRepository)
	deleteBookCommandHandler := book_application.NewDeleteBookCommandHandler(bookRepository)
	updateBookCommandHandler := book_application.NewUpdateBookCommandHandler(bookRepository)
	getAllBooksQueryHandler := book_application.NewGetAllBooksQueryHandler(bookRepository)
	getBookByIDQueryHandler := book_application.NewGetBookByIDQueryHandler(bookRepository)
	printBookOnBookCreatedEvent := book_application.NewPrintBookOnBookCreatedEventHandler()

	bookModuleServices := &BookModuleServices{
		CreateBookCommandHandler:           createBookCommandHandler,
		DeleteBookCommandHandler:           deleteBookCommandHandler,
		UpdateBookCommandHandler:           updateBookCommandHandler,
		GetAllBooksQueryHandler:            getAllBooksQueryHandler,
		GetBookByIDQueryHandler:            getBookByIDQueryHandler,
		PrintBookOnBookCreatedEventHandler: printBookOnBookCreatedEvent,
	}

	registerBookCommandHandlers(commonServices, bookModuleServices)
	registerBookQueryHandlers(commonServices, bookModuleServices)
	registerBookEventHandlers(commonServices, bookModuleServices)
	registerBookRoutes(commonServices, httpServices)

	return bookModuleServices
}

func registerBookCommandHandlers(commonServices *CommonServices, bookModuleServices *BookModuleServices) {
	registerCommandOrPanic(
		commonServices.CommandBus,
		&book_application.CreateBookCommand{},
		bookModuleServices.CreateBookCommandHandler,
	)
	registerCommandOrPanic(
		commonServices.CommandBus,
		&book_application.DeleteBookCommand{},
		bookModuleServices.DeleteBookCommandHandler,
	)
	registerCommandOrPanic(
		commonServices.CommandBus,
		&book_application.UpdateBookCommand{},
		bookModuleServices.UpdateBookCommandHandler,
	)
}

func registerBookQueryHandlers(commonServices *CommonServices, bookModuleServices *BookModuleServices) {
	registerQueryOrPanic(
		commonServices.QueryBus,
		&book_application.GetAllBooksQuery{},
		bookModuleServices.GetAllBooksQueryHandler,
	)
	registerQueryOrPanic(
		commonServices.QueryBus,
		&book_application.GetBookByIDQuery{},
		bookModuleServices.GetBookByIDQueryHandler,
	)
}

func registerBookEventHandlers(commonServices *CommonServices, bookModuleServices *BookModuleServices) {
	registerEvent(
		*commonServices.EventBus,
		book_domain.BookCreatedEventName,
		bookModuleServices.PrintBookOnBookCreatedEventHandler,
	)

}

func registerBookRoutes(commonServices *CommonServices, httpServices *HttpServices) {
	httpServices.Router.Post(
		"/api/books",
		book_http.NewPostBookController(
			commonServices.CommandBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)

	httpServices.Router.Delete(
		"/api/books/{book_id}",
		book_http.NewDeleteBookController(
			commonServices.CommandBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)

	httpServices.Router.Put(
		"/api/books/{book_id}",
		book_http.NewPutBookController(
			commonServices.CommandBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)

	httpServices.Router.Get(
		"/api/books",
		book_http.NewGetBooksController(
			commonServices.UlidProvider,
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)

	httpServices.Router.Get(
		"/api/books/{book_id}",
		book_http.NewGetBookController(
			commonServices.UlidProvider,
			commonServices.QueryBus,
			httpServices.JsonApiResponseMiddleware,
		),
	)
}
