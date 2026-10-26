package handler

import "github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"

type Notification struct {
	errors []validation.Error
}

func NewNotification() *Notification {
	return &Notification{errors: make([]validation.Error, 0)}
}

func (n *Notification) Append(err validation.Error) validation.Handler {
	n.errors = append(n.errors, err)
	return n
}

func (n *Notification) AppendHandler(handler validation.Handler) validation.Handler {
	n.errors = append(n.errors, handler.GetErrors()...)
	return n
}

func (n *Notification) GetErrors() []validation.Error {
	return n.errors
}

func (n *Notification) HasError() bool {
	return len(n.errors) > 0
}
