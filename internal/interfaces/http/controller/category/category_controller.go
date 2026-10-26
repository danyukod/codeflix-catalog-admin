package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
)

type Controller struct {
	create *createcategory.CategoryUseCase
	list   *listcategory.CategoryUseCase
	get    *getcategory.CategoryUseCase
	update *updatecategory.CategoryUseCase
	delete *deletecategory.CategoryUseCase
}

func NewController(
	create *createcategory.CategoryUseCase,
	list *listcategory.CategoryUseCase,
	get *getcategory.CategoryUseCase,
	update *updatecategory.CategoryUseCase,
	delete *deletecategory.CategoryUseCase,
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

	writeJSON(w, http.StatusCreated, responseFromCreate(*output))
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
		items = append(items, responseFromList(item))
	}
	writeJSON(w, http.StatusOK, categoryListResponse{CurrentPage: output.CurrentPage, PerPage: output.PerPage, Total: output.Total, Items: items})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var command getcategory.CategoryCommand
	output, err := c.get.Execute(r.Context(), command.With(id))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFromGet(*output))
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var request categoryRequest
	if !decodeJSON(w, r, &request) {
		return
	}
	var command updatecategory.CategoryCommand
	output, err := c.update.Execute(r.Context(), command.With(r.PathValue("id"), request.Name, request.Description, request.IsActive))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, responseFromUpdate(*output))
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
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
type categoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   *string `json:"deleted_at,omitempty"`
}

type categoryListResponse struct {
	CurrentPage int                `json:"current_page"`
	PerPage     int                `json:"per_page"`
	Total       int64              `json:"total"`
	Items       []categoryResponse `json:"items"`
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
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
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
	if errors.Is(err, strconv.ErrSyntax) || strings.Contains(err.Error(), "validation failed") || strings.Contains(err.Error(), "id cannot be empty") {
		status = http.StatusBadRequest
	}
	if strings.Contains(err.Error(), "was not found") {
		status = http.StatusNotFound
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func queryInt(r *http.Request, key string, fallback int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(key))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}
