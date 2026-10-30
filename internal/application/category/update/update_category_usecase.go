package update

import (
	"context"
	"fmt"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception/handler"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
)

type CategoryUseCase struct {
	gateway category.Gateway
}

func NewUpdateCategoryUsecase(gateway category.Gateway) *CategoryUseCase {
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

	aName := aCommand.name
	aDescription := aCommand.description
	isActive := aCommand.isActive

	aCategory, err := uc.gateway.FindById(ctx, anId)
	if err != nil {
		return nil, err
	}
	if aCategory == nil {
		return nil, &exception.NotFoundError{Entity: "category", ID: anId.GetValue()}
	}

	aCategory.Update(aName, aDescription, isActive)

	notification := handler.NewNotification()
	aCategory.Validate(notification)

	if notification.HasError() {
		return nil, validation.NewValidationError(notification.GetErrors())
	}

	err = uc.gateway.Update(ctx, aCategory)
	if err != nil {
		return nil, err
	}

	var updateCategoryOutput CategoryOutput
	out := updateCategoryOutput.From(*aCategory)

	return out, nil
}
