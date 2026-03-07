package update

import (
	"context"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateCategoryUsecase(t *testing.T) {
	gateway := new(mocks.CategoryGatewayMock)
	useCase := NewUpdateCategoryUsecase(gateway)

	t.Run("given a valid command when calls update category should return category id", func(t *testing.T) {
		ctx := context.TODO()
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

		gateway.On("FindById", expectedID).Return(aCategory, nil)
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
		})).Return(nil)

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.Equal(t, expectedID.GetValue(), actualOutput.GetId().GetValue())

		gateway.AssertExpectations(t)
	})
}
