package list

import (
	"context"
	"errors"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestListCategoryUseCase(t *testing.T) {
	ctx := context.TODO()

	t.Run("given a nil gateway when creates usecase should panic", func(t *testing.T) {
		assert.Panics(t, func() {
			NewListCategoryUseCase(nil)
		})
	})

	t.Run("given a valid command when calls list categories should return paginated output", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewListCategoryUseCase(gateway)

		cat1 := category.NewCategory("Filmes", "Categoria de filmes", true)
		cat2 := category.NewCategory("Series", "Categoria de series", true)

		expectedPage := pagination.New(1, 10, 2, []category.Category{*cat1, *cat2})

		var cmd CategoryCommand
		aCommand := cmd.With(1, 10, "", "name", "asc")

		gateway.On("FindAll", aCommand.ToSearchQuery()).Return(expectedPage, nil).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.Equal(t, 1, actualOutput.CurrentPage)
		assert.Equal(t, 10, actualOutput.PerPage)
		assert.Equal(t, int64(2), actualOutput.Total)
		assert.Len(t, actualOutput.Items, 2)
		assert.Equal(t, cat1.GetId().GetValue(), actualOutput.Items[0].GetId().GetValue())
		assert.Equal(t, cat1.GetName(), actualOutput.Items[0].GetName())
		assert.Equal(t, cat2.GetId().GetValue(), actualOutput.Items[1].GetId().GetValue())
		assert.Equal(t, cat2.GetName(), actualOutput.Items[1].GetName())

		gateway.AssertExpectations(t)
	})

	t.Run("given a nil command when calls list categories should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewListCategoryUseCase(gateway)

		actualOutput, err := useCase.Execute(ctx, nil)

		assert.Error(t, err)
		assert.EqualError(t, err, "aCommand cannot be nil")
		assert.Empty(t, actualOutput.Items)
		gateway.AssertNotCalled(t, "FindAll")
	})

	t.Run("given a valid command when find all returns an error should return error", func(t *testing.T) {
		gateway := new(mocks.CategoryGatewayMock)
		useCase := NewListCategoryUseCase(gateway)
		expectedErr := errors.New("generic gateway error")

		var cmd CategoryCommand
		aCommand := cmd.With(1, 10, "", "name", "asc")

		gateway.On("FindAll", aCommand.ToSearchQuery()).Return(nil, expectedErr).Once()

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Error(t, err)
		assert.EqualError(t, err, "generic gateway error")
		assert.Empty(t, actualOutput.Items)

		gateway.AssertExpectations(t)
	})
}
