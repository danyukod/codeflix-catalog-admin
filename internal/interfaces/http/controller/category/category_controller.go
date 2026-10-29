package category

import (
	"net/http"

	"github.com/danyukod/codeflix-catalog-admin/internal/application"
	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
	"github.com/danyukod/codeflix-catalog-admin/internal/domain/pagination"
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

func (c *Controller) Replace(w http.ResponseWriter, r *http.Request) {
	if !validID(w, r.PathValue("id")) {
		return
	}
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
	response := responseFromUpdate(*output)
	response.Links = categoryLinks(response.ID)
	writeJSON(w, http.StatusOK, response)
}

func (c *Controller) Patch(w http.ResponseWriter, r *http.Request) {
	if !validID(w, r.PathValue("id")) {
		return
	}
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
