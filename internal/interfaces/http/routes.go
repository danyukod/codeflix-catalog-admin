package http

import (
	stdhttp "net/http"

	"github.com/danyukod/codeflix-catalog-admin/docs"
	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
)

func NewRouter(controller *categorycontroller.Controller) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("/docs", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodGet {
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(docs.SwaggerUIHTML)
	})
	mux.HandleFunc("/openapi.yaml", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodGet {
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(docs.OpenAPIYAML)
	})
	mux.HandleFunc("/categories", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.URL.Path != "/categories" {
			stdhttp.NotFound(w, r)
			return
		}
		switch r.Method {
		case stdhttp.MethodGet:
			controller.List(w, r)
		case stdhttp.MethodPost:
			controller.Create(w, r)
		default:
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/categories/", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		switch r.Method {
		case stdhttp.MethodGet:
			controller.Get(w, r)
		case stdhttp.MethodPut, stdhttp.MethodPatch:
			controller.Update(w, r)
		case stdhttp.MethodDelete:
			controller.Delete(w, r)
		default:
			w.WriteHeader(stdhttp.StatusMethodNotAllowed)
		}
	})
	return mux
}
