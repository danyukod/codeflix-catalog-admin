package mocks

import (
	"context"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
	"github.com/stretchr/testify/mock"
)

type CategoryGatewayMock struct {
	mock.Mock
}

func (m *CategoryGatewayMock) Create(ctx context.Context, c *category.Category) (*category.Category, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	Category := args.Get(0)
	return Category.(*category.Category), args.Error(1)
}

func (m *CategoryGatewayMock) DeleteById(ctx context.Context, id category.Id) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindById(ctx context.Context, id category.Id) (*category.Category, error) {
	args := m.Called(ctx, id)

	var cat *category.Category
	if v := args.Get(0); v != nil {
		cat = v.(*category.Category)
	}

	return cat, args.Error(1)
}

func (m *CategoryGatewayMock) Update(ctx context.Context, c *category.Category) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindAll(ctx context.Context, query pagination.SearchQuery) (pagination.Pagination[category.Category], error) {
	args := m.Called(ctx, query)

	var page pagination.Pagination[category.Category]
	if v := args.Get(0); v != nil {
		page = v.(pagination.Pagination[category.Category])
	}

	return page, args.Error(1)
}

func (m *CategoryGatewayMock) ExistsByIDs(ctx context.Context, ids []category.Id) ([]category.Id, error) {
	args := m.Called(ctx, ids)

	var found []category.Id
	if v := args.Get(0); v != nil {
		found = v.([]category.Id)
	}

	return found, args.Error(1)
}
