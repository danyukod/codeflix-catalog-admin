package application

import "github.com/danyukod/codeflix-catalog-admin/internal/domain"

type CreateCategoryUseCase struct {
}

func (c *CreateCategoryUseCase) Execute() (*domain.Category, error) {
	return new(domain.Category), nil
}
