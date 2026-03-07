package update

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type CategoryOutput struct {
	id identifier.Identifier
}

func (u CategoryOutput) From(aCategory category.Category) *CategoryOutput {
	return &CategoryOutput{
		id: aCategory.GetId(),
	}
}

func (u CategoryOutput) GetId() identifier.Identifier { return u.id }
