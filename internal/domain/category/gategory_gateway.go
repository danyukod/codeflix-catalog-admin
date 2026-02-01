package category

import "github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"

type Gateway interface {
	Create(category *Category) error
	DeleteById(id Id) error
	FindById(id Id) (*Category, error)
	Update(category *Category) error
	FindAll(query pagination.SearchQuery) (pagination.Pagination[Category], error)
	ExistsByIDs(ids []Id) ([]Id, error)
}
