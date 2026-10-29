package category

import (
	"time"

	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
)

type categoryResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	DeletedAt   *string           `json:"deleted_at,omitempty"`
	Links       map[string]string `json:"_links"`
}

type categoryListResponse struct {
	CurrentPage int                `json:"current_page"`
	PerPage     int                `json:"per_page"`
	Total       int64              `json:"total"`
	Items       []categoryResponse `json:"items"`
	Links       map[string]string  `json:"_links"`
}

func responseFromList(o listcategory.CategoryOutput) categoryResponse {
	return categoryResponse{ID: o.GetId().GetValue(), Name: o.GetName(), Description: o.GetDescription(), IsActive: o.IsActive(), CreatedAt: o.GetCreatedAt().Format(time.RFC3339Nano), UpdatedAt: o.GetUpdatedAt().Format(time.RFC3339Nano)}
}

func responseFromGet(o getcategory.CategoryOutput) categoryResponse {
	return categoryResponse{ID: o.GetId().GetValue(), Name: o.GetName(), Description: o.GetDescription(), IsActive: o.IsActive(), CreatedAt: o.GetCreatedAt().Format(time.RFC3339Nano), UpdatedAt: o.GetUpdatedAt().Format(time.RFC3339Nano)}
}

func responseFromCreate(o createcategory.CategoryOutput) categoryResponse {
	return categoryResponse{ID: o.GetId().GetValue(), Name: o.GetName(), Description: o.GetDescription(), IsActive: o.IsActive(), CreatedAt: o.GetCreatedAt().Format(time.RFC3339Nano), UpdatedAt: o.GetUpdatedAt().Format(time.RFC3339Nano)}
}

func responseFromUpdate(o updatecategory.CategoryOutput) categoryResponse {
	return categoryResponse{ID: o.GetId().GetValue(), Name: o.GetName(), Description: o.GetDescription(), IsActive: o.IsActive(), CreatedAt: o.GetCreatedAt().Format(time.RFC3339Nano), UpdatedAt: o.GetUpdatedAt().Format(time.RFC3339Nano)}
}
