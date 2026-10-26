package category

import (
	"strings"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
)

const (
	NameMaxLength = 255
	NameMinLength = 3
)

type Validator struct {
	category *Category
	handler  validation.Handler
}

func NewCategoryValidator(aCategory *Category, aHandler validation.Handler) *Validator {
	return &Validator{
		category: aCategory,
		handler:  aHandler,
	}
}

func (v *Validator) Validate() {
	v.checkNameConstraints()
}

func (v *Validator) checkNameConstraints() {
	name := v.category.GetName()

	if strings.TrimSpace(name) == "" {
		v.handler.Append(validation.NewError("'name' should not be empty"))
		return
	}

	length := len(strings.TrimSpace(name))
	if length > NameMaxLength || length < NameMinLength {
		v.handler.Append(validation.NewError("'name' must be between 3 and 255 characters"))
	}
}
