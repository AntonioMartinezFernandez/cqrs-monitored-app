package bus

type InvalidDto struct {
	message string
}

func NewInvalidDto(message string) *InvalidDto {
	return &InvalidDto{message: message}
}

func (i InvalidDto) Error() string {
	return i.message
}
