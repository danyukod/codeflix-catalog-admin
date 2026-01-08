package application

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
)

type CreateCategoryUseCase struct {
}

func (c *CreateCategoryUseCase) Execute() (*category.Category, error) {
	return new(category.Category), nil
}
