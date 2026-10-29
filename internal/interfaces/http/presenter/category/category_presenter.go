package category

import (
	"time"

	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
)

type CategoryResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	DeletedAt   *string           `json:"deleted_at,omitempty"`
	Links       map[string]string `json:"_links"`
}

type CategoryListResponse struct {
	CurrentPage int                `json:"current_page"`
	PerPage     int                `json:"per_page"`
	Total       int64              `json:"total"`
	Items       []CategoryResponse `json:"items"`
	Links       map[string]string  `json:"_links"`
}

func PresentList(o listcategory.CategoryOutput) CategoryResponse {
	return presentCategory(o.GetId().GetValue(), o.GetName(), o.GetDescription(), o.IsActive(), o.GetCreatedAt(), o.GetUpdatedAt())
}

func PresentGet(o getcategory.CategoryOutput) CategoryResponse {
	return presentCategory(o.GetId().GetValue(), o.GetName(), o.GetDescription(), o.IsActive(), o.GetCreatedAt(), o.GetUpdatedAt())
}

func PresentCreate(o createcategory.CategoryOutput) CategoryResponse {
	return presentCategory(o.GetId().GetValue(), o.GetName(), o.GetDescription(), o.IsActive(), o.GetCreatedAt(), o.GetUpdatedAt())
}

func PresentUpdate(o updatecategory.CategoryOutput) CategoryResponse {
	return presentCategory(o.GetId().GetValue(), o.GetName(), o.GetDescription(), o.IsActive(), o.GetCreatedAt(), o.GetUpdatedAt())
}

func CategoryURL(id string) string { return "/categories/" + id }

func CategoryLinks(id string) map[string]string {
	return map[string]string{"self": CategoryURL(id), "collection": "/categories"}
}

func presentCategory(id, name, description string, active bool, createdAt, updatedAt time.Time) CategoryResponse {
	return CategoryResponse{
		ID: id, Name: name, Description: description, IsActive: active,
		CreatedAt: createdAt.Format(time.RFC3339Nano), UpdatedAt: updatedAt.Format(time.RFC3339Nano),
		Links: CategoryLinks(id),
	}
}
