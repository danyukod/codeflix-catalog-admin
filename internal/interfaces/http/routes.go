package http

import (
	stdhttp "net/http"

	"github.com/danyukod/codeflix-catalog-admin/docs"
	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
)

func NewRouter(controller *categorycontroller.Controller) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	mux.HandleFunc("GET /docs", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(docs.SwaggerUIHTML)
	})
	mux.HandleFunc("GET /openapi.yaml", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(docs.OpenAPIYAML)
	})
	mux.HandleFunc("GET /categories", controller.List)
	mux.HandleFunc("POST /categories", controller.Create)
	mux.HandleFunc("GET /categories/{id}", controller.Get)
	mux.HandleFunc("PUT /categories/{id}", controller.Update)
	mux.HandleFunc("PATCH /categories/{id}", controller.Update)
	mux.HandleFunc("DELETE /categories/{id}", controller.Delete)
	return mux
}
