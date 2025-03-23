package messaging

import "context"

type Opts func(messageConsumer *MessageConsumer) *MessageConsumer

func WithPoolSize(poolSize int) Opts {
	return func(messageConsumer *MessageConsumer) *MessageConsumer {
		if poolSize < 1 {
			messageConsumer.logger.Warn(
				context.Background(),
				"consumer pool size must be greater than 0, setting to default value of 10",
			)
			return messageConsumer
		}

		messageConsumer.poolSize = poolSize
		return messageConsumer
	}
}

func WithMessagesToConsume(messagesToConsume int) Opts {
	return func(messageConsumer *MessageConsumer) *MessageConsumer {
		if messagesToConsume < 1 {
			messageConsumer.logger.Warn(
				context.Background(),
				"messages to consume must be greater than 0, setting to default value of nil",
			)
			return messageConsumer
		}

		messageConsumer.messagesToConsume = &messagesToConsume
		return messageConsumer
	}
}
