package create

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception/handler"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewCreateCategoryUseCase(gateway category.Gateway) *CategoryUseCase {
	if gateway == nil {
		panic("gateway cannot be nil")
	}
	return &CategoryUseCase{gateway: gateway}
}

func (uc *CategoryUseCase) Execute(ctx context.Context, aCommand *CategoryCommand) (*CategoryOutput, error) {
	if aCommand == nil {
		return nil, fmt.Errorf("aCommand cannot be nil")
	}

	aName := aCommand.name
	aDescription := aCommand.description
	isActive := aCommand.isActive

	aCategory := category.NewCategory(aName, aDescription, isActive)

	notification := handler.NewNotification()
	aCategory.Validate(notification)
	if notification.HasError() {
		return nil, validation.NewValidationError(notification.GetErrors())
	}

	created, err := uc.gateway.Create(ctx, aCategory)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, fmt.Errorf("gateway returned nil category")
	}

	var createCategoryOutput CategoryOutput
	out := createCategoryOutput.From(*created)
	return out, nil
}
