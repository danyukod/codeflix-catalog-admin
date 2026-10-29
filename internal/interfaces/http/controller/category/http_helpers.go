package category

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/danyukod/codeflix-catalog-admin/internal/domain/validation"
	"github.com/google/uuid"
)

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
	if validationErr, ok := errors.AsType[*validation.ValidationError](err); ok {
		writeJSON(w, http.StatusBadRequest, map[string]any{"type": "https://api.codeflix.com/errors/validation", "title": "Validation failed", "status": http.StatusBadRequest, "detail": validationErr.Error(), "errors": validationErr.Messages})
		return
	}
	status := http.StatusInternalServerError
	if errors.Is(err, strconv.ErrSyntax) || strings.Contains(err.Error(), "id cannot be empty") || strings.Contains(err.Error(), "valid UUID") {
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
