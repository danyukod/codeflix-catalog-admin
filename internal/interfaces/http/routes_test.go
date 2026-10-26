package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDocumentationRoutes(t *testing.T) {
	router := NewRouter(nil)
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
