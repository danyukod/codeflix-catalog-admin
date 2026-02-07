package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception/handler"
)

type CreateCategoryUseCase struct {
	gateway category.Gateway
}

func NewCreateCategoryUseCase(gateway category.Gateway) *CreateCategoryUseCase {
	if gateway == nil {
		panic("gateway cannot be nil")
	}
	return &CreateCategoryUseCase{gateway: gateway}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, aCommand *CreateCategoryCommand) (*CreateCategoryOutput, error) {
	if aCommand == nil {
		return nil, fmt.Errorf("aCommand cannot be nil")
	}

	aName := aCommand.name
	aDescription := aCommand.description
	isActive := aCommand.isActive

	// Create domain entity from command
	aCategory := category.NewCategory(aName, aDescription, isActive)

	// Validate
	notification := handler.NewNotification()
	aCategory.Validate(notification)
	if notification.HasError() {
		return nil, fmt.Errorf("validation failed: %s", joinValidationMessages(notification))
	}

	// Persist
	created, err := uc.gateway.Create(aCategory)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, fmt.Errorf("gateway returned nil category")
	}

	var createCategoryOutput CreateCategoryOutput
	out := createCategoryOutput.From(*created)
	return out, nil
}

func joinValidationMessages(n *handler.Notification) string {
	errs := n.GetErrors()
	if len(errs) == 0 {
		return ""
	}

	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Message)
	}
	return strings.Join(msgs, "; ")
}
