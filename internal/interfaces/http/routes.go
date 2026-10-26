package http

import (
	stdhttp "net/http"

	categorycontroller "github.com/danyukod/codeflix-catalog-admin/internal/interfaces/http/controller/category"
	"github.com/swaggo/http-swagger"
)

func NewRouter(controller *categorycontroller.Controller) stdhttp.Handler {
	mux := stdhttp.NewServeMux()
	mux.Handle("/swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
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
