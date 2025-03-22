package command

import "github.com/AntonioMartinezFernandez/cqrs-monitored-app/pkg/bus"

type FailedCommand struct {
	command        bus.Dto
	handler        CommandHandler
	timesProcessed int
}

type CommandAlreadyRegistered struct {
	message     string
	commandName string
}

func (i CommandAlreadyRegistered) Error() string {
	return i.message
}

func NewCommandAlreadyRegistered(message string, commandName string) CommandAlreadyRegistered {
	return CommandAlreadyRegistered{message: message, commandName: commandName}
}

type CommandNotRegistered struct {
	message     string
	commandName string
}

func (i CommandNotRegistered) Error() string {
	return i.message
}

func NewCommandNotRegistered(message string, commandName string) CommandNotRegistered {
	return CommandNotRegistered{message: message, commandName: commandName}
}
