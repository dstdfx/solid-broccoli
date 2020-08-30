package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
)

const getSummeryQuery = `SELECT COUNT(1) FROM positions WHERE domain = $1`

// Position represents a single domain's position.
type Position struct {
	Domain   string     `db:"domain"`
	URL      string     `db:"keyword"`
	Position int        `db:"position"`
	Keyword  string     `db:"keyword"`
	Volume   int        `db:"volume"`
	Results  int        `db:"results"`
	CPC      float64    `db:"cpc"`
	Updated  *time.Time `db:"datetime"`
}

// DomainSummary represents a total number of positions for domain.
type DomainSummary struct {
	Domain         string
	PositionsCount int
}

// GetSummary returns a total number of positions for the given domain.
func (pr *PositionRepo) GetSummary(ctx context.Context, domain string) (*DomainSummary, error) {
	row := pr.conn.QueryRowxContext(ctx, getSummeryQuery, domain)
	err := row.Err()
	if err != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, err
		}
		pr.log.Error("failed to execute query", zap.Error(err))

		return nil, err
	}

	var positionsCount int
	if err := row.Scan(&positionsCount); err != nil {
		pr.log.Error("failed to scan positions count", zap.Error(err))

		return nil, err
	}

	return &DomainSummary{
		Domain:         domain,
		PositionsCount: positionsCount,
	}, nil
}
