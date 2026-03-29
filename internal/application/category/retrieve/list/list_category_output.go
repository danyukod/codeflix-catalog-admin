package list

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

func CategoryOutputFrom(aCategory category.Category) CategoryOutput {
	return CategoryOutput{
		id:          aCategory.GetId(),
		name:        aCategory.GetName(),
		description: aCategory.GetDescription(),
		isActive:    aCategory.IsActive(),
		createdAt:   aCategory.GetCreatedAt(),
		updatedAt:   aCategory.GetUpdatedAt(),
		deletedAt:   aCategory.GetDeletedAt(),
	}
}

func (o CategoryOutput) GetId() identifier.Identifier { return o.id }
func (o CategoryOutput) GetName() string              { return o.name }
func (o CategoryOutput) GetDescription() string       { return o.description }
func (o CategoryOutput) IsActive() bool               { return o.isActive }
func (o CategoryOutput) GetCreatedAt() time.Time      { return o.createdAt }
func (o CategoryOutput) GetUpdatedAt() time.Time      { return o.updatedAt }
func (o CategoryOutput) GetDeletedAt() *time.Time     { return o.deletedAt }
