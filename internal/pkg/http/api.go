package http

import (
	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	v1 "github.com/dstdfx/solid-broccoli/internal/pkg/http/v1"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

const (
	groupV1 = "/v1"
)

// InitAPIRouter configures HTTP router.
func InitAPIRouter(log *zap.Logger, b *backend.Backend) chi.Router {
	r := chi.NewRouter()
	r.Route(groupV1, func(r chi.Router) {
		r.Mount("/", v1.Routes(log, b))
	})

	return r
}
