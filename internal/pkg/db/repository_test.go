package db

import (
	"context"
	"testing"

	"github.com/dstdfx/solid-broccoli/internal/pkg/backend"
	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"github.com/dstdfx/solid-broccoli/internal/pkg/log"
	"github.com/dstdfx/solid-broccoli/internal/pkg/testutils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const dbSchema = `CREATE TABLE positions (keyword text, position integer, domain text, url text, volume integer, results integer, cpc float, updated datetime, primary key (domain, url, keyword));`

func prepareDBSchemaAndData(t *testing.T, conn *sqlx.DB) {
	_, err := conn.Exec(dbSchema)
	assert.NoError(t, err)
	_, err = conn.Exec(`
	insert into positions (keyword, position , domain, url, volume, results, cpc, updated) values 
("test1", 1, "ulmart.ru", "http://ulmart.ru/test1", 43, 40000, 1.22, 1495248847),
("test2", 2, "ulmart.ru", "http://ulmart.ru/test2", 55, 40000, 2.22, 1495248847),
("test3", 3, "ulmart.ru", "http://ulmart.ru/test3", 76, 40000, 3.22, 1495248847),
("test4", 7, "non-ulmart.ru", "http://nonulmart.ru/tests/2", 32, 40000, 5.22, 1495248847),
("test4", 11, "non-ulmart.ru", "http://nonulmart.ru/tests", 65, 40000, 5.22, 1495248847)`)
	assert.NoError(t, err)
}

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

	prepareDBSchemaAndData(t, b.DB)

	testDomain := "ulmart.ru"

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
