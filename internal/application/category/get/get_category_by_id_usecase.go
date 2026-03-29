package get

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewGetCategoryByIdUseCase(gateway category.Gateway) *CategoryUseCase {
	if gateway == nil {
		panic("gateway cannot be nil")
	}
	return &CategoryUseCase{gateway: gateway}
}

func (uc *CategoryUseCase) Execute(ctx context.Context, aCommand *CategoryCommand) (*CategoryOutput, error) {
	if aCommand == nil {
		return nil, fmt.Errorf("aCommand cannot be nil")
	}

	anId, err := category.FromString(aCommand.id)
	if err != nil {
		return nil, err
	}

	aCategory, err := uc.gateway.FindById(anId)
	if err != nil {
		return nil, err
	}
	if aCategory == nil {
		return nil, fmt.Errorf("category with ID %s was not found", anId.GetValue())
	}

	var output CategoryOutput
	return output.From(*aCategory), nil
}
