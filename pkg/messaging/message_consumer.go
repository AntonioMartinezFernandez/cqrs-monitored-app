package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/logger"
)

type MessageHandler interface {
	HandlerFunc() (HandlerFunc, error)
	Validator() (ValidatorFunc, error)
	MessageType() string
	MessageResolver() MessageResolver
}

type HandlerFunc func(ctx context.Context, msg Message) error
type ValidatorFunc func(ctx context.Context, msg *BaseMessage) error
type MessageResolver func(*BaseMessage) (Message, error)

type MessageConsumer struct {
	logger      logger.Logger
	retriever   MessageRetriever
	middlewares []Middleware

	poolSize          int
	messagesToConsume *int

	mut            *sync.Mutex
	subCtxCancelFn context.CancelFunc
	wg             *sync.WaitGroup

	handlers         map[string]HandlerFunc
	validators       map[string]ValidatorFunc
	messageResolvers map[string]MessageResolver
}

func NewMessageConsumer(
	logger logger.Logger,
	retriever MessageRetriever,
	middlewares []Middleware,
	options ...Opts,
) *MessageConsumer {
	messageConsumer := &MessageConsumer{
		logger:      logger,
		retriever:   retriever,
		middlewares: middlewares,

		poolSize:         1,
		mut:              &sync.Mutex{},
		handlers:         make(map[string]HandlerFunc),
		validators:       make(map[string]ValidatorFunc),
		messageResolvers: make(map[string]MessageResolver),
		wg:               &sync.WaitGroup{},

		messagesToConsume: nil,
	}

	for _, option := range options {
		messageConsumer = option(messageConsumer)
	}

	return messageConsumer
}

func (mc *MessageConsumer) RegisterHandler(messageHandler MessageHandler) {
	msgType := messageHandler.MessageType()
	messageResolver := messageHandler.MessageResolver()
	msgHandler, err := messageHandler.HandlerFunc()
	if err != nil {
		panic(err)
	}

	msgValidator, err := messageHandler.Validator()
	if err != nil {
		panic(err)
	}

	handlerWithMiddlewares := mc.applyMiddlewares(msgHandler)
	if mc.handlers[msgType] != nil {
		panic(fmt.Sprintf("message handler already registered for %s", msgType))
	}
	mc.handlers[msgType] = handlerWithMiddlewares

	if mc.validators[msgType] != nil {
		panic(fmt.Sprintf("message validator already registered for %s", msgType))
	}
	mc.validators[msgType] = msgValidator

	if mc.messageResolvers[msgType] != nil {
		panic(fmt.Sprintf("message instance already registered for %s", msgType))
	}
	mc.messageResolvers[msgType] = messageResolver
}

func (mc *MessageConsumer) Run(ctx context.Context) {
	// Define context for the goroutines that will handle messages
	subCtx, subCtxCancelFn := context.WithCancel(ctx)
	mc.subCtxCancelFn = subCtxCancelFn

	// Start retrieving messages from infra
	msgCh, err := mc.retriever.StartRetrieving(subCtx)
	if err != nil {
		mc.logger.Error(
			subCtx,
			"error retrieving messages from message queue",
			slog.String("error", err.Error()),
		)
		return
	}

	// Start pool of goroutines to handle messages
	mc.wg.Add(mc.poolSize)
	mc.startPool(subCtx, msgCh)
	mc.wg.Wait()
}

func (mc *MessageConsumer) Stop() {
	mc.subCtxCancelFn()
}

// startPool starts a pool of goroutines which listens to messages from the msgCh channel (each goroutine processes a single message)
func (mc *MessageConsumer) startPool(ctx context.Context, msgCh <-chan *BaseMessage) {
	msgCounter := 0

	// Start a goroutine to monitor the number of messages consumed if messagesToConsume is set
	if mc.messagesToConsume != nil {
		go mc.monitorMessageConsumption(ctx, &msgCounter)
	}

	for range mc.poolSize {
		go mc.worker(ctx, msgCh, &msgCounter)
	}
}

func (mc *MessageConsumer) monitorMessageConsumption(ctx context.Context, msgCounter *int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			mc.mut.Lock()
			if *msgCounter >= *mc.messagesToConsume {
				mc.mut.Unlock()
				mc.Stop()
				return
			}
			mc.mut.Unlock()
		}
	}
}

func (mc *MessageConsumer) worker(ctx context.Context, msgCh <-chan *BaseMessage, msgCounter *int) {
	defer mc.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case baseMessage, channelOpen := <-msgCh:
			if !channelOpen {
				return
			}

			// If messagesToConsume is set, increment the counter
			if mc.messagesToConsume != nil {
				mc.incrementMessageCounter(msgCounter)
			}

			// Process the message
			mc.processMessage(ctx, baseMessage)
		}
	}
}

func (mc *MessageConsumer) incrementMessageCounter(msgCounter *int) {
	mc.mut.Lock()
	*msgCounter++
	mc.mut.Unlock()
}

func (mc *MessageConsumer) processMessage(ctx context.Context, baseMessage *BaseMessage) {
	msgType := baseMessage.Type()

	// Get the handler
	handler := mc.handlers[msgType]
	if handler == nil {
		mc.logger.Error(
			ctx,
			"no handler registered for message",
			slog.String("message_type", msgType),
		)
		baseMessage.Ack()
		return
	}

	// Get the validator
	validator := mc.validators[msgType]
	if validator == nil {
		mc.logger.Error(
			ctx,
			"no validator registered for message",
			slog.String("message_type", msgType),
		)
		baseMessage.Ack()
		return
	}

	// Validate the message
	if err := validator(ctx, baseMessage); err != nil {
		mc.logger.Error(
			ctx,
			"error validating message",
			slog.String("error", err.Error()),
			slog.String("message_type", msgType),
			slog.Any("message", baseMessage.Data()),
		)
		baseMessage.Ack()
		return
	}

	// Get the message resolver
	msgResolver := mc.messageResolvers[msgType]
	if msgResolver == nil {
		mc.logger.Error(
			ctx,
			"no message resolver registered for message",
			slog.String("message_type", msgType),
		)
		baseMessage.Ack()
		return
	}

	// Resolve the message
	msg, err := msgResolver(baseMessage)
	if err != nil {
		mc.logger.Error(
			ctx,
			"error resolving message",
			slog.String("error", err.Error()),
			slog.String("message_type", msgType),
			slog.Any("message", baseMessage.Data()),
		)
		baseMessage.Ack()
		return
	}

	// Handle the message
	if err := handler(ctx, msg); err != nil {
		mc.logger.Error(
			ctx,
			"error handling message",
			slog.String("error", err.Error()),
		)
	}

	// Acknowledge the message
	baseMessage.Ack()
}

// applyMiddlewares applies middlewares to the handler
func (mc *MessageConsumer) applyMiddlewares(handler HandlerFunc) HandlerFunc {
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		handler = func(handler HandlerFunc, m Middleware) HandlerFunc {
			return func(ctx context.Context, msg Message) error {
				return m(ctx, msg, func(ctx context.Context, msg Message) error {
					return handler(ctx, msg)
				})
			}
		}(handler, mc.middlewares[i])
	}

	return handler
}
