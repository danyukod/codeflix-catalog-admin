package validation

type Error struct {
	Message string
}

func NewError(message string) Error {
	return Error{Message: message}
}

type ValidationHandler interface {
	Append(err Error) ValidationHandler
	AppendHandler(handler ValidationHandler) ValidationHandler
	GetErrors() []Error
	HasError() bool
}

func HasError(h ValidationHandler) bool {
	return len(h.GetErrors()) > 0
}
