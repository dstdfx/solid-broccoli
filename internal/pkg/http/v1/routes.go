package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/db"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

const (
	summaryURL   = "/summary"
	positionsURL = "/positions"

	defaultLimitPositionsPerPage = 10
)

// Routes initializes v1 handler.
func Routes(log *zap.Logger, b *backend.Backend) http.Handler {
	r := chi.NewRouter().
		With(chimiddleware.Recoverer).
		With(SetRequestID(log)).
		With(RequestLogger(log)).
		With(SetContextLogger(log)).
		With(RequireDomainName)

	// GET /v1/summary/<domain-name>
	r.Get(fmt.Sprintf("%s/{%s}", summaryURL, domainNameParam), summaryHandler(b))

	// GET /v1/positions/<domain-name>?orderBy=<field>&page=<page-num>
	r.Get(fmt.Sprintf("%s/{%s}", positionsURL, domainNameParam), positionsHandler(b))

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
		summary, err := repo.GetSummary(req.Context(), domain)
		if err != nil {
			log.Error("failed to get summary", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		JSON(w, NewSummaryResponse(domain, summary))
	}
}

func NewSummaryResponse(domain string, positions int) interface{} {
	return struct {
		Domain         string `json:"domain"`
		PositionsCount int    `json:"positions_count"`
	}{Domain: domain, PositionsCount: positions}
}

func positionsHandler(b *backend.Backend) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log, err := GetContextLogger(req.Context())
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)

			return
		}
		domain := GetDomainName(req.Context())

		// TODO: put query params logic into middlewares

		// Extract query params for page number and order by field
		var pageNum int
		rawPageNum := req.URL.Query().Get("page")
		if pageNum, err = strconv.Atoi(rawPageNum); err != nil {
			pageNum = 1
		}

		orderBy := req.URL.Query().Get("orderBy")
		if err := validateOrderByField(orderBy); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			JSON(w, map[string]string{"error": err.Error()})

			return
		}

		// Init repository and query domain's positions
		repo := db.NewPositionRepo(log, b.DB)
		positions, err := repo.GetPositions(
			req.Context(),
			domain,
			orderBy,
			defaultLimitPositionsPerPage,
			defaultLimitPositionsPerPage*(pageNum-1))
		if err != nil {
			log.Error("failed to get positions", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)

			return
		}

		// Write response
		w.WriteHeader(http.StatusOK)
		JSON(w, NewPositionsResponse(domain, positions))
	}
}

func NewPositionsResponse(domain string, positions []*db.Position) interface{} {
	return struct {
		Domain    string         `json:"domain"`
		Positions []*db.Position `json:"positions"`
	}{Domain: domain, Positions: positions}
}
