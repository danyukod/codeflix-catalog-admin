package category

import (
	"strings"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception/handler"
	"github.com/stretchr/testify/assert"
)

func TestNewCategory(t *testing.T) {
	t.Run(`Given Valid Params When Call New Category Then Instantiate a Category`, func(t *testing.T) {
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"

		expectedCategory := NewCategory(expectedName, expectedDescription, true)

		assert.NotNil(t, expectedCategory.GetId().GetIdentifier())
		assert.NotNil(t, expectedCategory)
		assert.Equal(t, expectedName, expectedCategory.GetName())
		assert.Equal(t, expectedDescription, expectedCategory.GetDescription())
		assert.True(t, expectedCategory.IsActive())
		assert.NotNil(t, expectedCategory.GetCreatedAt())
		assert.NotNil(t, expectedCategory.GetUpdatedAt())
		assert.Nil(t, expectedCategory.GetDeletedAt())
	})

	t.Run(`Given Valid Params When Call New Category Inactive Then Instantiate a Category With DeletedAt`, func(t *testing.T) {
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"

		category := NewCategory(expectedName, expectedDescription, false)

		assert.NotNil(t, category)
		assert.False(t, category.IsActive())
		assert.NotNil(t, category.GetDeletedAt())
	})
}

func TestCategory_Validate(t *testing.T) {
	t.Run("Given Invalid Empty Name When Call Validate Then Should Receive Error", func(t *testing.T) {
		category := NewCategory("  ", "Description", true)
		handler := handler.NewNotification()

		category.Validate(handler)

		assert.True(t, handler.HasError())
		assert.Equal(t, 1, len(handler.GetErrors()))
		assert.Equal(t, "'name' should not be empty", handler.GetErrors()[0].Message)
	})

	t.Run("Given Name Length Less Than 3 When Call Validate Then Should Receive Error", func(t *testing.T) {
		category := NewCategory("Fi", "Description", true)
		handler := handler.NewNotification()

		category.Validate(handler)

		assert.True(t, handler.HasError())
		assert.Equal(t, 1, len(handler.GetErrors()))
		assert.Equal(t, "'name' must be between 3 and 255 characters", handler.GetErrors()[0].Message)
	})

	t.Run("Given Name Length Greater Than 255 When Call Validate Then Should Receive Error", func(t *testing.T) {
		invalidName := strings.Repeat("a", 256)
		category := NewCategory(invalidName, "Description", true)
		handler := handler.NewNotification()

		category.Validate(handler)

		assert.True(t, handler.HasError())
		assert.Equal(t, 1, len(handler.GetErrors()))
	})

	t.Run("Given Valid Category When Call Validate Then Should Not Receive Error", func(t *testing.T) {
		category := NewCategory("Filmes", "Description", true)
		handler := handler.NewNotification()

		category.Validate(handler)

		assert.False(t, handler.HasError())
	})
}
