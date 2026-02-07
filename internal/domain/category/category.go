package category

import (
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/identifier"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
)

type Category struct {
	id          identifier.Identifier
	name        string
	description string
	active      bool
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   *time.Time
}

func NewCategory(aName string, aDescription string, isActive bool) *Category {
	id := Unique()
	now := time.Now().UTC()
	var deletedAt *time.Time
	if !isActive {
		deletedAt = &now
	}

	return &Category{
		id:          id,
		name:        aName,
		description: aDescription,
		active:      isActive,
		createdAt:   now,
		updatedAt:   now,
		deletedAt:   deletedAt,
	}
}

func (c *Category) Validate(handler validation.ValidationHandler) {
	NewCategoryValidator(c, handler).Validate()
}

func (c *Category) Deactivate() {
	c.active = false
	c.updatedAt = time.Now().UTC()
	c.deletedAt = &c.updatedAt
}

func (c *Category) Activate() {
	c.active = true
	c.updatedAt = time.Now().UTC()
	c.deletedAt = nil
}

func (c *Category) Update(aName, aDescription string, isActive bool) {
	if isActive {
		c.Activate()
	} else {
		c.Deactivate()
	}

	c.name = aName
	c.description = aDescription
	c.updatedAt = time.Now().UTC()
}

func (c *Category) GetId() identifier.Identifier { return c.id }
func (c *Category) GetName() string              { return c.name }
func (c *Category) GetDescription() string       { return c.description }
func (c *Category) IsActive() bool               { return c.active }
func (c *Category) GetCreatedAt() time.Time      { return c.createdAt }
func (c *Category) GetUpdatedAt() time.Time      { return c.updatedAt }
func (c *Category) GetDeletedAt() *time.Time     { return c.deletedAt }
