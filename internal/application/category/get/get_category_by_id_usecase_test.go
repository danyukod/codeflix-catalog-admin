package get

import (
	"context"
	"errors"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoryByIdUsecase(t *testing.T) {
	ctx := context.TODO()

	t.Run("given a nil gateway when creates usecase should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewGetCategoryByIdUseCase(nil)
		})
	})

	t.Run("given a valid id when calls get category by id should return category output", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewGetCategoryByIdUseCase(gateway)
		aCategory := category.NewCategory("Filmes", "A categoria mais assistida", true)
		expectedID := aCategory.GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.Equal(t, expectedID.GetValue(), actualOutput.GetId().GetValue())
		assert.Equal(t, aCategory.GetName(), actualOutput.GetName())
		assert.Equal(t, aCategory.GetDescription(), actualOutput.GetDescription())
		assert.Equal(t, aCategory.IsActive(), actualOutput.IsActive())
		assert.Equal(t, aCategory.GetCreatedAt(), actualOutput.GetCreatedAt())
		assert.Equal(t, aCategory.GetUpdatedAt(), actualOutput.GetUpdatedAt())
		assert.Equal(t, aCategory.GetDeletedAt(), actualOutput.GetDeletedAt())

		gateway.AssertExpectations(t)
	})

	t.Run("given a valid id when find by id returns an error should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewGetCategoryByIdUseCase(gateway)
		expectedID := category.NewCategory("Filmes", "", true).GetId()
		expectedErr := errors.New("generic gateway error")

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(nil, expectedErr).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "generic gateway error")
		assert.Nil(t, actualOutput)

		gateway.AssertExpectations(t)
	})

	t.Run("given a nil command when calls get category by id should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewGetCategoryByIdUseCase(gateway)

		actualOutput, err := useCase.Execute(ctx, nil)

		assert.Error(t, err)
		assert.EqualError(t, err, "aCommand cannot be nil")
		assert.Nil(t, actualOutput)
		gateway.AssertNotCalled(t, "FindById")
	})

	t.Run("given an invalid id when calls get category by id should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewGetCategoryByIdUseCase(gateway)

		var cmd CategoryCommand
		aCommand := cmd.With("")

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "id cannot be empty")
		assert.Nil(t, actualOutput)
		gateway.AssertNotCalled(t, "FindById")
	})

	t.Run("given a not found id when calls get category by id should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewGetCategoryByIdUseCase(gateway)
		expectedID := category.NewCategory("Filmes", "", true).GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(nil, nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "category with ID "+expectedID.GetValue()+" was not found")
		assert.Nil(t, actualOutput)

		gateway.AssertExpectations(t)
	})
}
