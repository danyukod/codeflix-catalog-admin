package create

import (
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type CategoryOutput struct {
	id          identifier.Identifier
	name        string
	description string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func (c CategoryOutput) From(aCategory category.Category) *CategoryOutput {
	return &CategoryOutput{
		id: aCategory.GetId(), name: aCategory.GetName(), description: aCategory.GetDescription(),
		isActive: aCategory.IsActive(), createdAt: aCategory.GetCreatedAt(), updatedAt: aCategory.GetUpdatedAt(), deletedAt: aCategory.GetDeletedAt(),
	}
}

func (c CategoryOutput) GetId() identifier.Identifier {
	return c.id
}
func (c CategoryOutput) GetName() string          { return c.name }
func (c CategoryOutput) GetDescription() string   { return c.description }
func (c CategoryOutput) IsActive() bool           { return c.isActive }
func (c CategoryOutput) GetCreatedAt() time.Time  { return c.createdAt }
func (c CategoryOutput) GetUpdatedAt() time.Time  { return c.updatedAt }
func (c CategoryOutput) GetDeletedAt() *time.Time { return c.deletedAt }
