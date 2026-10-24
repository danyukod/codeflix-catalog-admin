package validation

type Error struct {
	Message string
}

func NewError(message string) Error {
	return Error{Message: message}
}

type Handler interface {
	Append(err Error) Handler
	AppendHandler(handler Handler) Handler
	GetErrors() []Error
	HasError() bool
}
