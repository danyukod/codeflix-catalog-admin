package exception

type DomainException struct {
	errors  []error
	message string
}

func NewDomainException(errors []error, message string) *DomainException {
	return &DomainException{errors, message}
}

func (e *DomainException) Errors() []error {
	return e.Errors()
}
