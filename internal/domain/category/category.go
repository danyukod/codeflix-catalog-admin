package category

import (
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
)

type Category struct {
	id          identifier.Identifier
	name        string
	description string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func NewCategory(name, description string, isActive bool) *Category {
	return &Category{
		id:          Unique(),
		name:        name,
		description: description,
		isActive:    isActive,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

func (c *Category) GetName() string {
	return c.name
}

func (c *Category) GetDescription() string {
	return c.description
}

func (c *Category) IsActive() bool {
	return c.isActive
}

func (c *Category) GetCreatedAt() time.Time {
	return c.createdAt
}

func (c *Category) GetUpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Category) GetDeletedAt() *time.Time {
	return c.deletedAt
}
