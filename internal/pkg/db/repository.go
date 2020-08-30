package db

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// PositionRepo represents a data access layer to 'position' table.
type PositionRepo struct {
	conn *sqlx.DB
	log  *zap.Logger
}

// NewPositionRepo returns new instance of PositionRepo.
func NewPositionRepo(log *zap.Logger, conn *sqlx.DB) *PositionRepo {
	return &PositionRepo{
		conn: conn,
		log:  log,
	}
}
