package category

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
	"github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/shared"
)

func TestWriteErrorReturnsStructuredValidationProblem(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := validation.NewValidationError([]validation.Error{validation.NewError("'name' should not be empty")})

	shared.WriteError(recorder, err)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	body := recorder.Body.String()
	for _, expected := range []string{"Validation failed", "'name' should not be empty", "\"errors\""} {
		if !strings.Contains(body, expected) {
			t.Fatalf("response %q does not contain %q", body, expected)
		}
	}
}
