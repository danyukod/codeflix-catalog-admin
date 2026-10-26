package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	createcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/tests/mocks"
)

func TestDocumentationRoutes(t *testing.T) {
	router := NewRouter(testController())
	for _, tc := range []struct {
		path, contentType, contains string
	}{
		{"/docs", "text/html", "url: '/openapi.yaml'"},
		{"/openapi.yaml", "application/yaml", "openapi: 3.1.2"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if response.Code != http.StatusOK {
				t.Fatalf("GET %s: status %d", tc.path, response.Code)
			}
			if got := response.Header().Get("Content-Type"); !strings.HasPrefix(got, tc.contentType) {
				t.Fatalf("GET %s: Content-Type %q", tc.path, got)
			}
			if !strings.Contains(response.Body.String(), tc.contains) {
				t.Fatalf("GET %s: response missing %q", tc.path, tc.contains)
			}
		})
	}
}

func TestRouterRejectsUnsupportedRoutes(t *testing.T) {
	router := NewRouter(testController())
	for _, tc := range []struct {
		name, method, path string
		wantStatus         int
	}{
		{"wrong method for docs", http.MethodPost, "/docs", http.StatusMethodNotAllowed},
		{"wrong method for category", http.MethodPost, "/categories/123", http.StatusMethodNotAllowed},
		{"extra category path segment", http.MethodGet, "/categories/123/extra", http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, nil))
			if response.Code != tc.wantStatus {
				t.Fatalf("%s %s: status %d, want %d", tc.method, tc.path, response.Code, tc.wantStatus)
			}
		})
	}
}

func testController() *categorycontroller.Controller {
	gateway := &mocks.CategoryGatewayMock{}
	return categorycontroller.NewController(
		createcategory.NewCreateCategoryUseCase(gateway),
		listcategory.NewListCategoryUseCase(gateway),
		getcategory.NewGetCategoryByIdUseCase(gateway),
		updatecategory.NewUpdateCategoryUsecase(gateway),
		deletecategory.NewDeleteCategoryUseCase(gateway),
	)
}
