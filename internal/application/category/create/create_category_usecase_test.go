package create

import (
	"context"
	"errors"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategoryUseCase(t *testing.T) {
	gateway := new(mocks.CategoryGatewayMock)
	useCase := NewCreateCategoryUseCase(gateway)
	matchesCreatedCategory := func(expectedName, expectedDescription string, expectedIsActive bool) func(*category.Category) bool {
		return func(c *category.Category) bool {
			if c == nil {
				return false
			}

			if c.GetName() != expectedName || c.GetDescription() != expectedDescription || c.IsActive() != expectedIsActive {
				return false
			}

			if c.GetId() == nil || c.GetCreatedAt().IsZero() || c.GetUpdatedAt().IsZero() {
				return false
			}

			deletedAt := c.GetDeletedAt()
			deletedAtOK := (expectedIsActive && deletedAt == nil) || (!expectedIsActive && deletedAt != nil && !deletedAt.IsZero())
			return deletedAtOK
		}
	}

	t.Run("Given a valid command when calls create category should return category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedCategory := category.NewCategory(expectedName, expectedDescription, true)

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, true)

		gateway.
			On("Create", mock.MatchedBy(matchesCreatedCategory(expectedName, expectedDescription, true))).
			Return(expectedCategory, nil).
			Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.NotNil(t, actualOutput.GetId())
		gateway.AssertExpectations(t)
	})

	t.Run("Given a valida command with an invalid name when calls create category should return error", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := ""
		expectedDescription := "A categoria mais assistida"
		expectedError := "validation failed: 'name' should not be empty"

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, true)

		actualOutput, err := useCase.Execute(ctx, aCommand)

		gateway.AssertNotCalled(t, "Create")
		assert.EqualError(t, err, expectedError)
		assert.Nil(t, actualOutput)
	})

	t.Run("Given a valid command with a inactive category when calls create category should return a inactive category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedCategory := category.NewCategory(expectedName, expectedDescription, false)

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, false)

		gateway.
			On("Create", mock.MatchedBy(matchesCreatedCategory(expectedName, expectedDescription, false))).
			Return(expectedCategory, nil).
			Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.NotNil(t, actualOutput.GetId())
		gateway.AssertExpectations(t)
	})

	t.Run("Given a valid command when gateway returns error then should a error notification", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedError := "random error of gateway"

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, true)

		gateway.
			On("Create", mock.MatchedBy(matchesCreatedCategory(expectedName, expectedDescription, true))).
			Return(nil, errors.New(expectedError)).
			Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, actualOutput)
		assert.EqualError(t, err, expectedError)
		gateway.AssertExpectations(t)
	})

}
