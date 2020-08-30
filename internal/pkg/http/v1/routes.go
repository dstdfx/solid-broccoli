package v1

import (
	"net/http"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const (
	summaryURL   = "/summary"
	positionsURL = "/positions"
)

// Routes initializes v1 handler.
func Routes(log *zap.Logger, b *backend.Backend) http.Handler {
	r := chi.NewRouter()
	// TODO: add request-id middleware
	// TODO: add request logger

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
