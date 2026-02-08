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

func TestCreateCategoryUsecase(t *testing.T) {
	gateway := new(mocks.CategoryGatewayMock)
	useCase := NewCreateCategoryUseCase(gateway)

	t.Run("Given a valid command when calls create category should return category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedCategory := category.NewCategory(expectedName, expectedDescription, expectedIsActive)

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, expectedIsActive)

		gateway.
			On("Create", mock.MatchedBy(func(c *category.Category) bool {
				if c == nil {
					return false
				}
				return c.GetName() == expectedName &&
					c.GetDescription() == expectedDescription &&
					c.IsActive() == expectedIsActive &&
					c.GetId() != nil &&
					!c.GetCreatedAt().IsZero() &&
					!c.GetUpdatedAt().IsZero() &&
					c.GetDeletedAt() == nil
			})).
			Return(expectedCategory, nil).
			Once()

		// run the code-under-test here...

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.NotNil(t, actualOutput.GetId())
		gateway.AssertExpectations(t) // validates "once" + predicate matched
	})

	t.Run("Given a valida command with an invalid name when calls create category should return error", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := ""
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedError := "validation failed: 'name' should not be empty"

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, expectedIsActive)

		actualOutput, err := useCase.Execute(ctx, aCommand)

		gateway.AssertNotCalled(t, "Create")
		assert.EqualError(t, err, expectedError)
		assert.Nil(t, actualOutput)
	})

	t.Run("Given a valid command with a inactive category when calls create category should return a inactive category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := false
		expectedCategory := category.NewCategory(expectedName, expectedDescription, expectedIsActive)

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, expectedIsActive)

		gateway.
			On("Create", mock.MatchedBy(func(c *category.Category) bool {
				if c == nil {
					return false
				}
				return c.GetName() == expectedName &&
					c.GetDescription() == expectedDescription &&
					c.IsActive() == expectedIsActive &&
					c.GetId() != nil &&
					!c.GetCreatedAt().IsZero() &&
					!c.GetUpdatedAt().IsZero() &&
					!c.GetDeletedAt().IsZero()
			})).
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
		expectedIsActive := true
		expectedError := "random error of gateway"

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, expectedIsActive)

		gateway.
			On("Create", mock.MatchedBy(func(c *category.Category) bool {
				if c == nil {
					return false
				}
				return c.GetName() == expectedName &&
					c.GetDescription() == expectedDescription &&
					c.IsActive() == expectedIsActive &&
					c.GetId() != nil &&
					!c.GetCreatedAt().IsZero() &&
					!c.GetUpdatedAt().IsZero() &&
					c.GetDeletedAt() == nil
			})).
			Return(nil, errors.New(expectedError)).
			Once()

		// run the code-under-test here...

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, actualOutput)
		assert.EqualError(t, err, expectedError)
		gateway.AssertExpectations(t) // validates "once" + predicate matched
	})
}
