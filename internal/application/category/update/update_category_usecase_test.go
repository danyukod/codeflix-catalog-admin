package update

import (
	"context"
	"errors"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCategoryUsecase(t *testing.T) {
	gateway := new(mocks.CategoryGatewayMock)
	useCase := NewUpdateCategoryUsecase(gateway)
	ctx := context.TODO()

	t.Run("given a valid command when calls update category should return category id", func(t *testing.T) {
		aCategory := category.NewCategory("Film", "", true)

		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedID := aCategory.GetId()

		createdAt := aCategory.GetCreatedAt()
		updatedAt := aCategory.GetUpdatedAt()

		var cmd CategoryCommand
		aCommand := cmd.With(
			expectedID.GetValue(),
			expectedName,
			expectedDescription,
			expectedIsActive,
		)

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()
		gateway.On("Update", mock.MatchedBy(func(updatedCategory *category.Category) bool {
			if updatedCategory == nil {
				return false
			}

			return updatedCategory.GetId().GetValue() == expectedID.GetValue() &&
				updatedCategory.GetName() == expectedName &&
				updatedCategory.GetDescription() == expectedDescription &&
				updatedCategory.IsActive() == expectedIsActive &&
				updatedCategory.GetCreatedAt().Equal(createdAt) &&
				updatedCategory.GetUpdatedAt().After(updatedAt) &&
				updatedCategory.GetDeletedAt() == nil
		})).Return(nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.Equal(t, expectedID.GetValue(), actualOutput.GetId().GetValue())

		gateway.AssertExpectations(t)
	})

	t.Run("given a valid command when gateway returns a generic error should return error", func(t *testing.T) {
		aCategory := category.NewCategory("Film", "", true)

		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedID := aCategory.GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(
			expectedID.GetValue(),
			expectedName,
			expectedDescription,
			expectedIsActive,
		)

		expectedErr := errors.New("generic gateway error")

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()
		gateway.On("Update", mock.AnythingOfType("*category.Category")).Return(expectedErr).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "generic gateway error")
		assert.Nil(t, actualOutput)

		gateway.AssertExpectations(t)
	})

	t.Run("given a invalid command when calls update category should return error", func(t *testing.T) {
		aCategory := category.NewCategory("Film", "", true)
		expectedName := ""
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedID := aCategory.GetId()

		createdAt := aCategory.GetCreatedAt()
		updatedAt := aCategory.GetUpdatedAt()

		var cmd CategoryCommand
		aCommand := cmd.With(
			expectedID.GetValue(),
			expectedName,
			expectedDescription,
			expectedIsActive,
		)

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()
		gateway.On("Update", mock.MatchedBy(func(updatedCategory *category.Category) bool {
			if updatedCategory == nil {
				return false
			}

			return updatedCategory.GetId().GetValue() == expectedID.GetValue() &&
				updatedCategory.GetName() == expectedName &&
				updatedCategory.GetDescription() == expectedDescription &&
				updatedCategory.IsActive() == expectedIsActive &&
				updatedCategory.GetCreatedAt().Equal(createdAt) &&
				updatedCategory.GetUpdatedAt().After(updatedAt) &&
				updatedCategory.GetDeletedAt() != nil
		})).Return(nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.Equal(t, expectedID.GetValue(), actualOutput.GetId().GetValue())

		gateway.AssertExpectations(t)
	})

	t.Run("given a invalid command when calls update category should return error", func(t *testing.T) {
		aCategory := category.NewCategory("Film", "", true)
		expectedName := ""
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedID := aCategory.GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(
			expectedID.GetValue(),
			expectedName,
			expectedDescription,
			expectedIsActive,
		)

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		gateway.AssertCalled(t, "FindById", expectedID)
		gateway.AssertNotCalled(t, "Update")
		assert.EqualError(t, err, "validation failed: 'name' should not be empty")
		assert.Nil(t, actualOutput)
	})

}
