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

	// 1. Teste do caminho feliz.
	t.Run("Given a valid command when calls update category should return category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedCategory := category.NewCategory("Film", "", true)

		expectedId := expectedCategory.GetId().GetValue()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true

		var cmd UpdateCategoryCommand
		aCommand := cmd.With(expectedId, expectedName, expectedDescription, expectedIsActive)

		gateway.On("FindById", expectedId).Return(expectedCategory, nil)
		gateway.On("Update", mock.Anything).Return(nil)

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		gateway.AssertExpectations(t)
		gateway.AssertNumberOfCalls(t, "FindById", 1)
		gateway.AssertNumberOfCalls(t, "Update", 1)
		assert.Equal(t, expectedId, actualOutput.id.GetValue())
		assert.Equal(t, expectedName, actualOutput.name)
		assert.Equal(t, expectedDescription, actualOutput.description)
		assert.Equal(t, true, actualOutput.isActive)
	})
	// 2. Teste passando uma propriedade invalida (name)
	// 3. Teste atualizando uma categoria para inativa.
	// 4. Teste simulando um erro generico vindo do gateway.
}
