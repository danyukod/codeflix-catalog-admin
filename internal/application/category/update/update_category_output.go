package update

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

func (u CategoryOutput) From(aCategory category.Category) *CategoryOutput {
	return &CategoryOutput{
		id: aCategory.GetId(), name: aCategory.GetName(), description: aCategory.GetDescription(),
		isActive: aCategory.IsActive(), createdAt: aCategory.GetCreatedAt(), updatedAt: aCategory.GetUpdatedAt(), deletedAt: aCategory.GetDeletedAt(),
	}
}

func (u CategoryOutput) GetId() identifier.Identifier { return u.id }
func (u CategoryOutput) GetName() string              { return u.name }
func (u CategoryOutput) GetDescription() string       { return u.description }
func (u CategoryOutput) IsActive() bool               { return u.isActive }
func (u CategoryOutput) GetCreatedAt() time.Time      { return u.createdAt }
func (u CategoryOutput) GetUpdatedAt() time.Time      { return u.updatedAt }
func (u CategoryOutput) GetDeletedAt() *time.Time     { return u.deletedAt }
