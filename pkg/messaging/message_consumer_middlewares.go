package messaging

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"
)

type Middleware func(ctx context.Context, msg Message, next HandlerFunc) error

// PanicRecovererMiddleware recovers from panics and returns an error.
func PanicRecovererMiddleware() Middleware {
	return func(ctx context.Context, msg Message, next HandlerFunc) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.WithStack(PanicRecoveredError{V: r, Stacktrace: string(debug.Stack())})
			}
		}()
		return next(ctx, msg)
	}
}

// PanicRecoveredError is an error that is returned when a panic occurs in a handler.
type PanicRecoveredError struct {
	V          any
	Stacktrace string
}

func (p PanicRecoveredError) Error() string {
	return fmt.Sprintf("panic occurred: %#v, stacktrace: \n%s", p.V, p.Stacktrace)
}
