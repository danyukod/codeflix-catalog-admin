package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/exception/handler"
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

	aCategory, err := uc.gateway.FindById(anId)
	if err != nil {
		return nil, err
	}
	if aCategory == nil {
		return nil, fmt.Errorf("category with ID %s was not found", anId.GetValue())
	}

	aCategory.Update(aName, aDescription, isActive)

	notification := handler.NewNotification()
	aCategory.Validate(notification)

	if notification.HasError() {
		return nil, fmt.Errorf("validation failed: %s", joinValidationMessages(notification))
	}

	err = uc.gateway.Update(aCategory)
	if err != nil {
		return nil, err
	}

	var updateCategoryOutput CategoryOutput
	out := updateCategoryOutput.From(*aCategory)

	return out, nil
}

func joinValidationMessages(n *handler.Notification) string {
	errs := n.GetErrors()
	if len(errs) == 0 {
		return ""
	}

	messages := make([]string, 0, len(errs))
	for _, e := range errs {
		messages = append(messages, e.Message)
	}
	return strings.Join(messages, "; ")
}
