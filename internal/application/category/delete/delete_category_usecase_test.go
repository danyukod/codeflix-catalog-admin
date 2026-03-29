package delete

import (
	"context"
	"errors"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCategoryUsecase(t *testing.T) {
	ctx := context.TODO()

	t.Run("given a valid command when calls delete category should return nil error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewDeleteCategoryUseCase(gateway)
		aCategory := category.NewCategory("Film", "", true)
		expectedID := aCategory.GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()
		gateway.On("DeleteById", expectedID).Return(nil).Once()

		err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)

		gateway.AssertCalled(t, "FindById", expectedID)
		gateway.AssertCalled(t, "DeleteById", expectedID)
		gateway.AssertExpectations(t)
	})

	t.Run("given a valid command when gateway delete returns a generic error should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewDeleteCategoryUseCase(gateway)
		aCategory := category.NewCategory("Film", "", true)
		expectedID := aCategory.GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		expectedErr := errors.New("generic gateway error")

		gateway.On("FindById", expectedID).Return(aCategory, nil).Once()
		gateway.On("DeleteById", expectedID).Return(expectedErr).Once()

		err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "generic gateway error")

		gateway.AssertCalled(t, "FindById", expectedID)
		gateway.AssertCalled(t, "DeleteById", expectedID)
		gateway.AssertExpectations(t)
	})

	t.Run("given a valid command when find by id returns an error should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewDeleteCategoryUseCase(gateway)
		expectedID := category.NewCategory("Film", "", true).GetId()
		expectedErr := errors.New("generic gateway error")

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(nil, expectedErr).Once()

		err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "generic gateway error")
		gateway.AssertCalled(t, "FindById", expectedID)
		gateway.AssertNotCalled(t, "DeleteById")
		gateway.AssertExpectations(t)
	})

	t.Run("given a not found id when calls delete category should return nil error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewDeleteCategoryUseCase(gateway)
		expectedID := category.NewCategory("Film", "", true).GetId()

		var cmd CategoryCommand
		aCommand := cmd.With(expectedID.GetValue())

		gateway.On("FindById", expectedID).Return(nil, nil).Once()

		err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		gateway.AssertCalled(t, "FindById", expectedID)
		gateway.AssertNotCalled(t, "DeleteById")
		gateway.AssertExpectations(t)
	})
}
