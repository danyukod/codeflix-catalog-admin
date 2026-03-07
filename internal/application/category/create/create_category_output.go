package create

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type CategoryOutput struct {
	id identifier.Identifier
}

func (c CategoryOutput) From(aCategory category.Category) *CategoryOutput {
	return &CategoryOutput{
		id: aCategory.GetId(),
	}
}

func (c CategoryOutput) GetId() identifier.Identifier {
	return c.id
}
