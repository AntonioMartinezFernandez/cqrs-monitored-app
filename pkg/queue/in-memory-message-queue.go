package queue

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
)

var _ MessageQueue = (*InMemoryMessageQueue)(nil)

type InMemoryMessageQueue struct {
	mut    *sync.RWMutex
	logger logger.Logger

	messages               []Message
	handlers               map[string][]MessageHandler
	maxConcurrentProcesses int

	started      bool
	subCtxCancel context.CancelFunc

	msgChan     chan Message
	processesWG sync.WaitGroup
}

func NewInMemoryMessageQueue(logger logger.Logger, maxConcurrentProcesses int) *InMemoryMessageQueue {
	return &InMemoryMessageQueue{
		mut:                    &sync.RWMutex{},
		logger:                 logger,
		messages:               []Message{},
		handlers:               make(map[string][]MessageHandler),
		maxConcurrentProcesses: maxConcurrentProcesses,

		started:      false,
		subCtxCancel: nil,

		msgChan: make(chan Message, maxConcurrentProcesses),
	}
}

func (mq *InMemoryMessageQueue) Publish(message Message) error {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	reflectedMessageName, err := mq.reflectedMessageName(message)
	if err != nil {
		return err
	}

	if _, ok := mq.handlers[*reflectedMessageName]; !ok {
		return fmt.Errorf("no message handlers registered for message type %s", message.MessageType())
	}

	mq.messages = append(mq.messages, message)

	return nil
}

func (mq *InMemoryMessageQueue) RegisterHandler(message Message, handler MessageHandler) error {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	reflectedMessageName, err := mq.reflectedMessageName(message)
	if err != nil {
		return err
	}

	for _, h := range mq.handlers[*reflectedMessageName] {
		if h.Name() == handler.Name() {
			return fmt.Errorf("handler %s already registered for message type %s", handler.Name(), *reflectedMessageName)
		}
	}

	mq.handlers[*reflectedMessageName] = append(mq.handlers[*reflectedMessageName], handler)

	return nil
}

func (mq *InMemoryMessageQueue) DeregisterHandler(message Message, handler MessageHandler) error {
	mq.mut.Lock()
	defer mq.mut.Unlock()

	reflectedMessageName, err := mq.reflectedMessageName(message)
	if err != nil {
		return err
	}

	handlers, ok := mq.handlers[*reflectedMessageName]
	if !ok {
		return nil
	}

	for i, h := range handlers {
		if h.Name() == handler.Name() {
			mq.handlers[*reflectedMessageName] = slices.Delete(handlers, i, i+1)
			break
		}
	}

	return nil
}

func (mq *InMemoryMessageQueue) Start(ctx context.Context) error {
	if mq.started {
		return fmt.Errorf("message queue already started")
	}

	mq.started = true

	subCtx, cancel := context.WithCancel(ctx)
	mq.subCtxCancel = cancel

	go mq.startMessageProcessor(subCtx)

	return nil
}

func (mq *InMemoryMessageQueue) EnqueuedMessages() int {
	mq.mut.RLock()
	defer mq.mut.RUnlock()

	return len(mq.messages)
}

func (mq *InMemoryMessageQueue) startMessageProcessor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			mq.logger.Info(ctx, "stopping message queue because context was cancelled", slog.Int("pending_messages", mq.EnqueuedMessages()))
			return
		default:
			msgsLen := mq.EnqueuedMessages()
			if msgsLen > 0 {
				// Retrieve first maxConcurrentProcesses messages from the queue and send them to the message channel
				elementsToProcess := min(msgsLen, mq.maxConcurrentProcesses)
				for i := 0; i < elementsToProcess && i < len(mq.messages); i++ {
					mq.msgChan <- mq.messages[i]
				}

				// Remove the first maxConcurrentProcesses messages from the queue
				mq.mut.Lock()
				mq.messages = slices.Delete(mq.messages, 0, elementsToProcess)
				mq.mut.Unlock()

				// Process messages concurrently
				for range elementsToProcess {
					mq.processesWG.Add(1)
					go func() {
						defer mq.processesWG.Done()
						mq.processMessage(ctx, <-mq.msgChan)
					}()
				}
				mq.processesWG.Wait()
			} else {
				// If there are no messages to process, wait for a while before checking again
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (mq *InMemoryMessageQueue) processMessage(ctx context.Context, msg Message) {
	reflectedMessageName, err := mq.reflectedMessageName(msg)
	if err != nil {
		mq.logger.Error(ctx, "error processing message from message queue", slog.String("error", err.Error()))
		return
	}
	// Get all handlers for the message type
	handlers, ok := mq.handlers[*reflectedMessageName]
	if !ok {
		mq.logger.Error(ctx, "no handlers registered for message type", slog.String("message_type", *reflectedMessageName))
		return
	}
	if len(handlers) == 0 {
		mq.logger.Error(ctx, "no handlers registered for message type", slog.String("message_type", *reflectedMessageName))
		return
	}

	// Process the message with all handlers concurrently
	for _, handler := range handlers {
		wg := sync.WaitGroup{}
		wg.Add(1)

		startTime := time.Now()

		go func(handler MessageHandler) {
			defer wg.Done()

			executions := 0

			handlerCtx, cancel := context.WithTimeout(ctx, msg.MaxTimeout())
			defer cancel()

			for executions < msg.MaxRetries() {
				err := handler.Handle(handlerCtx, msg)

				if err == nil {
					mq.logger.Debug(
						ctx,
						"message processed",
						slog.String("message_type", *reflectedMessageName),
						slog.String("handler", handler.Name()),
						slog.Duration("duration_ms", time.Duration(time.Since(startTime).Milliseconds())),
					)
					break
				}

				mq.logger.Error(
					ctx,
					"error processing message",
					slog.String("message_type", *reflectedMessageName),
					slog.String("handler", handler.Name()),
					slog.String("error", err.Error()),
				)

				executions++
			}
		}(handler)

		wg.Wait()
	}
}

func (mq *InMemoryMessageQueue) reflectedMessageName(cmd any) (*string, error) {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("message must be a pointer to a struct")
	}

	name := value.String()

	return &name, nil
}
