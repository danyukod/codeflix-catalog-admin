package update

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewUpdateCategoryUsecase(gateway *mocks.CategoryGatewayMock) *CategoryUseCase {
	return &CategoryUseCase{
		gateway: gateway,
	}
}

func (uc *CategoryUseCase) Execute(ctx context.Context, aCommand *UpdateCategoryCommand) (*UpdateCategoryOutput, error) {
	if aCommand == nil {
		return nil, fmt.Errorf("aCommand cannot be nil")
	}

	aCategoryId, err := category.FromString(aCommand.id)
	if err != nil {
		return nil, err
	}

	aName := aCommand.name
	aDescription := aCommand.description
	isActive := aCommand.isActive

	aCategory, err := uc.gateway.FindById(aCategoryId)
	if err != nil {
		return nil, err
	}

	if aCategory == nil {
		return nil, fmt.Errorf("category not found")
	}

	updatedCategory, err := uc.gateway.Update()
}
