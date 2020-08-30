package v1

import (
	"net/http"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	summaryURL   = "/summary"
	positionsURL = "/positions"
)

// Routes initializes v1 handler.
func Routes(log *zap.Logger, b *backend.Backend) http.Handler {
	r := chi.NewRouter().
		With(middleware.Recoverer).
		With(SetRequestID(log)).
		With(RequestLogger(log)).
		With(SetContextLogger(log))
	r.Get(summaryURL, summaryHandler(b))
	r.Get(positionsURL, positionsHandler(b))

	return r
}

func summaryHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: implement me
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func positionsHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: implement me
		w.WriteHeader(http.StatusNotImplemented)
	}
}
