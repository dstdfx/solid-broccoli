package backend

import (
	"fmt"

	"github.com/dstdfx/solid-broccoli/internal/pkg/config"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

// Backend contains application connections to different external services and additional
// parameters that should be passed to API middlewares.
type Backend struct {
	Log *zap.Logger
	DB  *sqlx.DB
}

// New init new Backend instance.
func (b *Backend) New(log *zap.Logger) (*Backend, error) {
	// Init DB connection
	conn, err := sqlx.Connect("sqlite3", config.Config.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init DB connection: %w", err)
	}

	return &Backend{
		Log: log,
		DB:  conn,
	}, nil
}

// Shutdown method closes all backend connections.
func (b *Backend) Shutdown() {
	b.Log.Debug("backend shutdown")

	// Close DB connection
	if b.DB != nil {
		if err := b.DB.Close(); err != nil {
			b.Log.Warn("failed to close DB connection")
		}
	}
}
