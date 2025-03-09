package book_application

import (
	"context"
	"fmt"

	book_domain "github.com/AntonioMartinezFernandez/cqrs-monitored-app/internal/book/domain"
	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"
)

type PrintBookOnBookCreatedEventHandler struct{}

func NewPrintBookOnBookCreatedEventHandler() *PrintBookOnBookCreatedEventHandler {
	return &PrintBookOnBookCreatedEventHandler{}
}

func (eh *PrintBookOnBookCreatedEventHandler) Handle(ctx context.Context, event bus.Event) {
	evt, ok := event.(*book_domain.BookCreatedEvent)
	if !ok {
		fmt.Println("Invalid event")
	}

	fmt.Println("*** Received book created EVENT ***", evt.Data())
}
