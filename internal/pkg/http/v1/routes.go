package v1

import (
	"fmt"
	"net/http"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/db"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	summaryURL   = "/summary"
	positionsURL = "/positions"
)

// Routes initializes v1 handler.
func Routes(log *zap.Logger, b *backend.Backend) http.Handler {
	r := chi.NewRouter().
		With(chimiddleware.Recoverer).
		With(SetRequestID(log)).
		With(RequestLogger(log)).
		With(SetContextLogger(log))

	r.With(RequireDomainName).
		Get(fmt.Sprintf("%s/{%s}", summaryURL, domainNameParam), summaryHandler(b))

	// TODO: fixme
	r.Get(positionsURL, positionsHandler(b))

	return r
}

func summaryHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log, err := GetContextLogger(req.Context())
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}
		domain := GetDomainName(req.Context())

		repo := db.NewPositionRepo(log, b.DB)
		domainSummary, err := repo.GetSummary(req.Context(), domain)
		if err != nil {
			log.Error("failed to get summary", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		// TODO: write swagger models
		w.WriteHeader(http.StatusOK)
		JSON(w, newSummaryResponse(domainSummary))
	}
}

func newSummaryResponse(ds *db.DomainSummary) interface{} {
	return struct {
		Domain         string `json:"domain"`
		PositionsCount int    `json:"positions_count"`
	}{ds.Domain, ds.PositionsCount}
}

func positionsHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// TODO: implement me
		w.WriteHeader(http.StatusNotImplemented)
	}
}
