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

		expectedCategory := NewCategory(expectedName, expectedDescription, true)

		assert.NotNil(t, expectedCategory.id.GetIdentifier())
		assert.NotNil(t, expectedCategory)
		assert.Equal(t, expectedName, expectedCategory.GetName())
		assert.NotNil(t, expectedDescription, expectedCategory.GetDescription())
		assert.True(t, true, expectedCategory.IsActive())
		assert.NotNil(t, expectedCategory.GetCreatedAt())
		assert.NotNil(t, expectedCategory.GetUpdatedAt())
		assert.Nil(t, expectedCategory.GetDeletedAt())

	})

	t.Run(`Given Invalid Empty Name When Call New Category and Validate Then Should Receive Error`, func(t *testing.T) {
		expectedName := ""
		expectedDescription := "A categoria mais assistida"
		expectedErrorMessage := "name cannot be empty"
		expectedErrorCount := 1

		expectedCategory := NewCategory(expectedName, expectedDescription, true)
		err := expectedCategory.Validate()

		assert.NotNil(t, err)
		assert.Equal(t, expectedErrorMessage, err[0].Error())
		assert.Equal(t, expectedErrorCount, len(err))
	})

}
