package category

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ... existing code ...
func TestNewCategory(t *testing.T) {

	t.Run(`Given Valid Params When Call New Category Then Instantiate a Category`, func(t *testing.T) {
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true

		expectedCategory := NewCategory(expectedName, expectedDescription, expectedIsActive)

		assert.NotNil(t, expectedCategory)
		assert.NotNil(t, expectedCategory.GetId())
		assert.Equal(t, expectedName, expectedCategory.GetName())
		assert.NotNil(t, expectedDescription, expectedCategory.GetDescription())
		assert.True(t, true, expectedCategory.IsActive())
		assert.NotNil(t, expectedCategory.GetCreatedAt())
		assert.NotNil(t, expectedCategory.GetUpdatedAt())
		assert.Nil(t, expectedCategory.GetDeletedAt())

	})
}
