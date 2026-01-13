package category

import "github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"

type Validator struct {
	category *Category
	handler  validation.Handler
}

func NewValidator(category *Category, handler validation.Handler) *Validator {
	return &Validator{category, handler}
}

func (v *Validator) Validate() {

}
