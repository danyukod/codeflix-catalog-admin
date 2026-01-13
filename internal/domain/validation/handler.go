package validation

type Handler interface {
	Handle(error) error
	AppendError(error)
	AppendErrors([]error)
	Validate(Validation)
	HasErrors() bool
	Errors() []error
}

type Validation interface {
	Validate()
}
