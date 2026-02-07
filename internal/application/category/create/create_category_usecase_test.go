package create

import (
	"context"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCategoryUsecase(t *testing.T) {
	// 1. Teste do caminho feliz
	// 2. Teste passando uma propriedade invalida (name)
	// 3. Teste Criando uma categoria inativa
	// 4. Teste simulando um erro generico vindo do gateway

	t.Run("Given a valid command when calls create category should return category id", func(t *testing.T) {
		ctx := context.TODO()
		expectedName := "Filmes"
		expectedDescription := "A categoria mais assistida"
		expectedIsActive := true
		expectedCategory := category.NewCategory(expectedName, expectedDescription, expectedIsActive)

		var cmd CreateCategoryCommand
		aCommand := cmd.With(expectedName, expectedDescription, expectedIsActive)

		gateway := new(mocks.CategoryGatewayMock)

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

		useCase := NewCreateCategoryUseCase(gateway)

		actualOutput, err := useCase.Execute(ctx, aCommand)

		assert.Nil(t, err)
		assert.NotNil(t, actualOutput)
		assert.NotNil(t, actualOutput.GetId())
		gateway.AssertExpectations(t) // validates "once" + predicate matched
	})
}
