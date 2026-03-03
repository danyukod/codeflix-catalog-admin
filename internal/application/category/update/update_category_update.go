package update

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type UpdateCategoryOutput struct {
	id          identifier.Identifier
	name        string
	description string
	isActive    bool
}

func (u UpdateCategoryOutput) From(aCategory category.Category) *UpdateCategoryOutput {
	return &UpdateCategoryOutput{
		id:          aCategory.GetId(),
		name:        aCategory.GetName(),
		description: aCategory.GetDescription(),
		isActive:    aCategory.IsActive(),
	}
}

func (u UpdateCategoryOutput) GetId() identifier.Identifier { return u.id }
func (u UpdateCategoryOutput) GetName() string              { return u.name }
func (u UpdateCategoryOutput) GetDescription() string       { return u.description }
func (u UpdateCategoryOutput) IsActive() bool               { return u.isActive }
