package mocks

import (
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
	"github.com/stretchr/testify/mock"
)

type CategoryGatewayMock struct {
	mock.Mock
}

func (m *CategoryGatewayMock) Create(c *category.Category) (*category.Category, error) {
	args := m.Called(c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	Category := args.Get(0)
	return Category.(*category.Category), args.Error(1)
}

func (m *CategoryGatewayMock) DeleteById(id category.Id) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindById(id category.Id) (*category.Category, error) {
	args := m.Called(id)

	var cat *category.Category
	if v := args.Get(0); v != nil {
		cat = v.(*category.Category)
	}

	return cat, args.Error(1)
}

func (m *CategoryGatewayMock) Update(c *category.Category) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *CategoryGatewayMock) FindAll(query pagination.SearchQuery) (pagination.Pagination[category.Category], error) {
	args := m.Called(query)

	var page pagination.Pagination[category.Category]
	if v := args.Get(0); v != nil {
		page = v.(pagination.Pagination[category.Category])
	}

	return page, args.Error(1)
}

func (m *CategoryGatewayMock) ExistsByIDs(ids []category.Id) ([]category.Id, error) {
	args := m.Called(ids)

	var found []category.Id
	if v := args.Get(0); v != nil {
		found = v.([]category.Id)
	}

	return found, args.Error(1)
}
