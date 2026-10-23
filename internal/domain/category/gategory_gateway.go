package category

import (
	"context"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
)

type Gateway interface {
	Create(ctx context.Context, category *Category) (*Category, error)
	DeleteById(ctx context.Context, id Id) error
	FindById(ctx context.Context, id Id) (*Category, error)
	Update(ctx context.Context, category *Category) error
	FindAll(ctx context.Context, query pagination.SearchQuery) (pagination.Pagination[Category], error)
	ExistsByIDs(ctx context.Context, ids []Id) ([]Id, error)
}
