package category

import (
	"strings"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
)

const (
	NameMaxLength = 255
	NameMinLength = 3
)

type CategoryValidator struct {
	category *Category
	handler  validation.ValidationHandler
}

func NewCategoryValidator(aCategory *Category, aHandler validation.ValidationHandler) *CategoryValidator {
	return &CategoryValidator{
		category: aCategory,
		handler:  aHandler,
	}
}

func (v *CategoryValidator) Validate() {
	v.checkNameConstraints()
}

func (v *CategoryValidator) checkNameConstraints() {
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
