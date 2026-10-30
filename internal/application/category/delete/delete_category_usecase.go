package delete

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewDeleteCategoryUseCase(gateway category.Gateway) *CategoryUseCase {
	if gateway == nil {
		panic("gateway cannot be nil")
	}
	return &CategoryUseCase{gateway: gateway}
}

func (uc *CategoryUseCase) Execute(ctx context.Context, aCommand *CategoryCommand) error {
	if aCommand == nil {
		return fmt.Errorf("aCommand cannot be nil")
	}

	anId, err := category.FromString(aCommand.id)
	if err != nil {
		return err
	}

	aCategory, err := uc.gateway.FindById(ctx, anId)
	if err != nil {
		return err
	}
	if aCategory == nil {
		return &exception.NotFoundError{Entity: "category", ID: anId.GetValue()}
	}

	return uc.gateway.DeleteById(ctx, anId)
}
