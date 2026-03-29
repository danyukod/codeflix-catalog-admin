package list

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewListCategoryUseCase(gateway category.Gateway) *CategoryUseCase {
	if gateway == nil {
		panic("gateway cannot be nil")
	}
	return &CategoryUseCase{gateway: gateway}
}

func (uc *CategoryUseCase) Execute(ctx context.Context, aCommand *CategoryCommand) (pagination.Pagination[CategoryOutput], error) {
	if aCommand == nil {
		return pagination.Pagination[CategoryOutput]{}, fmt.Errorf("aCommand cannot be nil")
	}

	result, err := uc.gateway.FindAll(aCommand.ToSearchQuery())
	if err != nil {
		return pagination.Pagination[CategoryOutput]{}, err
	}

	return pagination.Map(result, CategoryOutputFrom), nil
}
