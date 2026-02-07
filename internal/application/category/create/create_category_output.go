package create

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type CreateCategoryOutput struct {
	id identifier.Identifier
}

func (c CreateCategoryOutput) From(aCategory category.Category) *CreateCategoryOutput {
	return &CreateCategoryOutput{
		id: aCategory.GetId(),
	}
}

func (c CreateCategoryOutput) GetId() identifier.Identifier {
	return c.id
}
