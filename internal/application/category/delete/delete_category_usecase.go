package delete

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
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

	aCategory, err := uc.gateway.FindById(anId)
	if err != nil {
		return err
	}
	if aCategory == nil {
		return nil
	}

	return uc.gateway.DeleteById(anId)
}
