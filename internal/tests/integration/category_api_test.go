//go:build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/application/category/create"
	deletecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/delete"
	getcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/get"
	listcategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/retrieve/list"
	updatecategory "github.com/danyukod/codeflix-catalog-admin/internal/application/category/update"
	"github.com/danyukod/codeflix-catalog-admin/internal/infrastructure/category"
	"github.com/danyukod/codeflix-catalog-admin/internal/infrastructure/database"
	apphttp "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http"
	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type testApp struct {
	db     *sql.DB
	server *httptest.Server
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()

	ctx := context.Background()
	container, err := mysqlcontainer.Run(
		ctx,
		"mysql:8.4",
		mysqlcontainer.WithDatabase("shop_test"),
		mysqlcontainer.WithUsername("shop_test"),
		mysqlcontainer.WithPassword("shop_test"),
	)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, testcontainers.TerminateContainer(container))
	})

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "3306/tcp")
	require.NoError(t, err)
	portNumber, err := strconv.Atoi(port.Port())
	require.NoError(t, err)

	dbConnection, err := database.NewConnection(ctx, database.Config{
		Host:               host,
		Port:               portNumber,
		Name:               "shop_test",
		User:               "shop_test",
		Password:           "shop_test",
		SSLMode:            "disable",
		MaxOpenConnections: 5,
		MaxIdleConnections: 5,
	})
	require.NoError(t, err)
	t.Cleanup(dbConnection.Close)

	require.NoError(t, dbConnection.RunMigrations(ctx))

	repo := category.NewRepository(dbConnection.Pool())
	controller := categorycontroller.NewController(
		create.NewCreateCategoryUseCase(repo),
		listcategory.NewListCategoryUseCase(repo),
		getcategory.NewGetCategoryByIdUseCase(repo),
		updatecategory.NewUpdateCategoryUsecase(repo),
		deletecategory.NewDeleteCategoryUseCase(repo),
	)
	server := httptest.NewServer(apphttp.NewRouter(controller))
	t.Cleanup(server.Close)

	return &testApp{db: dbConnection.Pool(), server: server}
}

func TestCategoryAPI(t *testing.T) {
	app := newTestApp(t)

	var created struct {
		ID string `json:"id"`
	}
	response := app.request(t, http.MethodPost, "/categories", `{"name":"Films","description":"Feature films","is_active":true}`)
	require.Equal(t, http.StatusCreated, response.StatusCode)
	require.NotEmpty(t, response.Header.Get("Location"))
	require.NoError(t, json.NewDecoder(response.Body).Decode(&created))
	require.NotEmpty(t, created.ID)
	require.NoError(t, response.Body.Close())

	response = app.request(t, http.MethodGet, "/categories/"+created.ID, "")
	require.Equal(t, http.StatusOK, response.StatusCode)
	var found map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&found))
	require.Equal(t, "Films", found["name"])
	require.Equal(t, "/categories/"+created.ID, found["_links"].(map[string]any)["self"])
	require.NoError(t, response.Body.Close())

	response = app.request(t, http.MethodPut, "/categories/"+created.ID, `{"name":"Updated films","description":"Updated description","is_active":true}`)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.NoError(t, response.Body.Close())

	response = app.request(t, http.MethodGet, "/categories", "")
	require.Equal(t, http.StatusOK, response.StatusCode)
	var listed struct {
		Total int64 `json:"total"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&listed))
	require.Equal(t, int64(1), listed.Total)
	require.NoError(t, response.Body.Close())

	response = app.request(t, http.MethodPatch, "/categories/"+created.ID, `{"description":"Patched description"}`)
	require.Equal(t, http.StatusOK, response.StatusCode)
	var patched map[string]any
	require.NoError(t, json.NewDecoder(response.Body).Decode(&patched))
	require.Equal(t, "Updated films", patched["name"])
	require.Equal(t, "Patched description", patched["description"])
	require.True(t, patched["is_active"].(bool))
	require.NoError(t, response.Body.Close())

	response = app.request(t, http.MethodDelete, "/categories/"+created.ID, "")
	require.Equal(t, http.StatusNoContent, response.StatusCode)
	require.NoError(t, response.Body.Close())

	var active bool
	err := app.db.QueryRowContext(t.Context(), "SELECT active FROM categories WHERE id = ?", created.ID).Scan(&active)
	require.NoError(t, err)
	require.False(t, active)
}

func TestCategoryAPIValidationError(t *testing.T) {
	app := newTestApp(t)
	response := app.request(t, http.MethodPost, "/categories", `{"name":"","description":"Invalid","is_active":true}`)
	require.Equal(t, http.StatusBadRequest, response.StatusCode)
	var problem struct {
		Type   string   `json:"type"`
		Title  string   `json:"title"`
		Status int      `json:"status"`
		Errors []string `json:"errors"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&problem))
	require.Equal(t, "https://api.codeflix.com/errors/validation", problem.Type)
	require.Equal(t, "Validation failed", problem.Title)
	require.Equal(t, http.StatusBadRequest, problem.Status)
	require.Contains(t, problem.Errors, "'name' should not be empty")
	require.NoError(t, response.Body.Close())
}

func TestCategoryAPINotFoundByID(t *testing.T) {
	app := newTestApp(t)
	notFoundID := "00000000-0000-0000-0000-000000000000"

	for _, test := range []struct {
		method string
		path   string
		body   string
	}{
		{http.MethodGet, "/categories/" + notFoundID, ""},
		{http.MethodPut, "/categories/" + notFoundID, `{"name":"Films","description":"Feature films","is_active":true}`},
		{http.MethodPatch, "/categories/" + notFoundID, `{"description":"Feature films"}`},
		{http.MethodDelete, "/categories/" + notFoundID, ""},
	} {
		response := app.request(t, test.method, test.path, test.body)
		require.Equal(t, http.StatusNotFound, response.StatusCode, test.method)
		require.NoError(t, response.Body.Close())
	}
}

func (app *testApp) request(t *testing.T, method, path, body string) *http.Response {
	t.Helper()

	request, err := http.NewRequestWithContext(t.Context(), method, app.server.URL+path, bytes.NewBufferString(body))
	require.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.server.Client().Do(request)
	require.NoError(t, err)
	return response
}
