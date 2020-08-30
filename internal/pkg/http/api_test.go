package http

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"github.com/dstdfx/solid-broccoli/internal/pkg/db"
	v1 "github.com/dstdfx/solid-broccoli/internal/pkg/http/v1"
	"github.com/dstdfx/solid-broccoli/internal/pkg/log"
	"github.com/dstdfx/solid-broccoli/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

// Tests for GET /v1/summary/<domain-name>

func TestGetSummaryOK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend.
	b, err := backend.New(logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	// Setup handlers
	router := InitAPIRouter(logger, b)

	// Test a request.
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/summary/%s", testutils.TestDomain)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			v1.NewSummaryResponse(testutils.TestDomain, 3),
		), w.Body.String())
}

// Tests for GET /v1/positions/<domain-name>

func TestGetPositionsOK(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend.
	b, err := backend.New(logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	// Setup handlers
	router := InitAPIRouter(logger, b)

	// Test a request.
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/positions/%s", testutils.TestDomain)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, r)

	repo := db.NewPositionRepo(logger, b.DB)
	positions, err := repo.GetPositions(context.Background(), testutils.TestDomain, "", 10, 0)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			v1.NewPositionsResponse(testutils.TestDomain, positions),
		), w.Body.String())
}

func TestGetPositions_BadOrderByField(t *testing.T) {
	// Check acceptance test flag
	if !testutils.IsAccTestEnabled(t) {
		return
	}

	// Init global app configuration
	testutils.InitTestConfig()

	// Initialize logger
	logger, err := log.InitLogger(log.InitLoggerOpts{
		Debug:     config.Config.Log.Debug,
		UseStdout: config.Config.Log.UseStdout,
		File:      config.Config.Log.File,
	})
	assert.NoError(t, err)

	// Prepare backend.
	b, err := backend.New(logger)
	assert.NoError(t, err)
	assert.NotEmpty(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	// Setup handlers
	router := InitAPIRouter(logger, b)

	// Test a request.
	w := httptest.NewRecorder()
	url := fmt.Sprintf("/v1/positions/%s?orderBy=wwwwat", testutils.TestDomain)
	r, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t,
		testutils.RespToJSON(t,
			map[string]string{"error": "positions can't be ordered by 'wwwwat' field"},
		), w.Body.String())
}
