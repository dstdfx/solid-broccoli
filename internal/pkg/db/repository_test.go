package db

import (
	"context"
	"testing"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"github.com/dstdfx/solid-broccoli/internal/pkg/log"
	"github.com/dstdfx/solid-broccoli/internal/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

const testDomain = "ulmart.ru"

func TestGetSummary(t *testing.T) {
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

	b, err := backend.New(logger)
	defer b.Shutdown()
	assert.NoError(t, err)
	assert.NotNil(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	expectedSummary := &DomainSummary{
		Domain:         testDomain,
		PositionsCount: 3,
	}

	repo := NewPositionRepo(logger, b.DB)
	got, err := repo.GetSummary(context.Background(), testDomain)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, expectedSummary, got)
}

func TestGetPositions_DefaultOrder(t *testing.T) {
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

	b, err := backend.New(logger)
	defer b.Shutdown()
	assert.NoError(t, err)
	assert.NotNil(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	repo := NewPositionRepo(logger, b.DB)
	// Query positions with default order by "volume"
	got, err := repo.GetPositions(context.Background(), testDomain, "", 1, 0)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)
	// Check that it has the lowest volume which is 43
	assert.Equal(t, 43, got[0].Volume)

	// Query next chunk with offset
	got, err = repo.GetPositions(context.Background(), testDomain, "", 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

	// Query the last chunk with offset
	got, err = repo.GetPositions(context.Background(), testDomain, "", 1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

	// Check that nothing is left after offset 3
	got, err = repo.GetPositions(context.Background(), testDomain, "", 1, 3)
	assert.NoError(t, err)
	assert.Empty(t, got)

	// TODO: use `require` instead of `assert`
}

func TestGetPositions_OrderByCPC(t *testing.T) {
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

	b, err := backend.New(logger)
	defer b.Shutdown()
	assert.NoError(t, err)
	assert.NotNil(t, b)

	testutils.PrepareDB(t, b.DB)
	defer testutils.TeardownDB(t, b.DB)

	repo := NewPositionRepo(logger, b.DB)
	// Query positions with default order by "volume"
	got, err := repo.GetPositions(context.Background(), testDomain, "cpc", 1, 0)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

	// Query next chunk with offset
	got, err = repo.GetPositions(context.Background(), testDomain, "cpc", 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

	// Query the last chunk with offset
	got, err = repo.GetPositions(context.Background(), testDomain, "cpc", 1, 2)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Len(t, got, 1)

	// Check that nothing is left after offset 3
	got, err = repo.GetPositions(context.Background(), testDomain, "cpc", 1, 3)
	assert.NoError(t, err)
	assert.Empty(t, got)
}
