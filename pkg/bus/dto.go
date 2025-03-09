package bus

type Dto interface {
	Type() string
}

type BlockOperationCommand interface {
	Dto
	BlockingKey() string
}
