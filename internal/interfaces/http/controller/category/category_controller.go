package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/danyukod/codeflix-catalog-admin/internal/application"
	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
	"github.com/google/uuid"
)

type Controller struct {
	create application.UseCase[*createcategory.CategoryCommand, *createcategory.CategoryOutput]
	list   application.UseCase[*listcategory.CategoryCommand, pagination.Pagination[listcategory.CategoryOutput]]
	get    application.UseCase[*getcategory.CategoryCommand, *getcategory.CategoryOutput]
	update application.UseCase[*updatecategory.CategoryCommand, *updatecategory.CategoryOutput]
	delete application.UnitUseCase[*deletecategory.CategoryCommand]
}

func NewController(
	create application.UseCase[*createcategory.CategoryCommand, *createcategory.CategoryOutput],
	list application.UseCase[*listcategory.CategoryCommand, pagination.Pagination[listcategory.CategoryOutput]],
	get application.UseCase[*getcategory.CategoryCommand, *getcategory.CategoryOutput],
	update application.UseCase[*updatecategory.CategoryCommand, *updatecategory.CategoryOutput],
	delete application.UnitUseCase[*deletecategory.CategoryCommand],
) *Controller {
	return &Controller{create: create, list: list, get: get, update: update, delete: delete}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var request categoryRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	var command createcategory.CategoryCommand
	output, err := c.create.Execute(r.Context(), command.With(request.Name, request.Description, request.IsActive))
	if err != nil {
		writeError(w, err)
		return
	}

	response := responseFromCreate(*output)
	response.Links = categoryLinks(response.ID)
	w.Header().Set("Location", categoryURL(response.ID))
	writeJSON(w, http.StatusCreated, response)
}

func (c *Controller) List(w http.ResponseWriter, r *http.Request) {
	page := queryInt(r, "page", 1)
	perPage := queryInt(r, "per_page", 10)
	var command listcategory.CategoryCommand
	output, err := c.list.Execute(r.Context(), command.With(page, perPage, r.URL.Query().Get("terms"), r.URL.Query().Get("sort"), r.URL.Query().Get("direction")))
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]categoryResponse, 0, len(output.Items))
	for _, item := range output.Items {
		response := responseFromList(item)
		response.Links = categoryLinks(response.ID)
		items = append(items, response)
	}
	writeJSON(w, http.StatusOK, categoryListResponse{CurrentPage: output.CurrentPage, PerPage: output.PerPage, Total: output.Total, Items: items, Links: map[string]string{"self": r.URL.RequestURI(), "collection": "/categories"}})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !validID(w, id) {
		return
	}
	var command getcategory.CategoryCommand
	output, err := c.get.Execute(r.Context(), command.With(id))
	if err != nil {
		writeError(w, err)
		return
	}
	response := responseFromGet(*output)
	response.Links = categoryLinks(response.ID)
	writeJSON(w, http.StatusOK, response)
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	if !validID(w, r.PathValue("id")) {
		return
	}
	if r.Method == http.MethodPatch {
		var request categoryPatchRequest
		if !decodeJSON(w, r, &request) {
			return
		}
		current, err := c.get.Execute(r.Context(), getcategory.CategoryCommand{}.With(r.PathValue("id")))
		if err != nil {
			writeError(w, err)
			return
		}
		name, description, active := current.GetName(), current.GetDescription(), current.IsActive()
		if request.Name != nil {
			name = *request.Name
		}
		if request.Description != nil {
			description = *request.Description
		}
		if request.IsActive != nil {
			active = *request.IsActive
		}
		var command updatecategory.CategoryCommand
		output, err := c.update.Execute(r.Context(), command.With(r.PathValue("id"), name, description, active))
		if err != nil {
			writeError(w, err)
			return
		}
		response := responseFromUpdate(*output)
		response.Links = categoryLinks(response.ID)
		writeJSON(w, http.StatusOK, response)
		return
	}
	var request categoryRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	if request.Name == "" {
		writeError(w, errors.New("PUT requires name"))
		return
	}
	var command updatecategory.CategoryCommand
	output, err := c.update.Execute(r.Context(), command.With(r.PathValue("id"), request.Name, request.Description, request.IsActive))
	if err != nil {
		writeError(w, err)
		return
	}
	response := responseFromUpdate(*output)
	response.Links = categoryLinks(response.ID)
	writeJSON(w, http.StatusOK, response)
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	if !validID(w, r.PathValue("id")) {
		return
	}
	var command deletecategory.CategoryCommand
	if err := c.delete.Execute(r.Context(), command.With(r.PathValue("id"))); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}
type categoryPatchRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}
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
	return categoryResponse{ID: o.GetId().GetValue(), Name: o.GetName(), Description: o.GetDescription(), IsActive: o.IsActive(), CreatedAt: o.GetCreatedAt().Format("2006-01-02T15:04:05.999Z07:00"), UpdatedAt: o.GetUpdatedAt().Format("2006-01-02T15:04:05.999Z07:00")}
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

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return false
	}
	if decoder.Decode(&struct{}{}) == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "request body must contain a single JSON value"})
		return false
	}
	return true
}
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, strconv.ErrSyntax) || strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "id cannot be empty") || strings.Contains(err.Error(), "valid UUID") || strings.Contains(err.Error(), "requires name") {
		status = http.StatusBadRequest
	}
	if strings.Contains(err.Error(), "was not found") {
		status = http.StatusNotFound
	}
	writeJSON(w, status, map[string]any{"type": "about:blank", "title": http.StatusText(status), "status": status, "detail": err.Error()})
}
func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func validID(w http.ResponseWriter, id string) bool {
	if _, err := uuid.Parse(id); err != nil {
		writeError(w, errors.New("id must be a valid UUID"))
		return false
	}
	return true
}

func categoryURL(id string) string { return "/categories/" + id }

func categoryLinks(id string) map[string]string {
	return map[string]string{"self": categoryURL(id), "collection": "/categories"}
}
