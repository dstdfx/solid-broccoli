package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

const (
	getSummaryQuery = `SELECT COUNT(1) FROM positions WHERE domain = $1`

	getPositionsQuery = `SELECT
				keyword,
				position,
				url,
				volume,
				results,
				cpc,
				date(updated, 'unixepoch')
		FROM positions
		WHERE domain = $1
		ORDER BY %s ASC
		LIMIT $2 OFFSET $3
`
)

// Position represents a single domain's position.
type Position struct {
	URL      string  `db:"url" json:"url"`
	Position int     `db:"position" json:"position"`
	Keyword  string  `db:"keyword" json:"keyword"`
	Volume   int     `db:"volume" json:"volume"`
	Results  int     `db:"results" json:"results"`
	CPC      float64 `db:"cpc" json:"cpc"`
	Updated  string  `db:"updated" json:"updated"`
}

// DomainSummary represents a total number of positions for domain.
type DomainSummary struct {
	Domain         string
	PositionsCount int
}

// GetSummary returns a total number of positions for the given domain.
func (pr *PositionRepo) GetSummary(ctx context.Context, domain string) (int, error) {
	row := pr.conn.QueryRowxContext(ctx, getSummaryQuery, domain)
	err := row.Err()
	if err != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return -1, err
		}
		pr.log.Error("failed to execute query", zap.Error(err))

		return -1, err
	}

	var positionsCount int
	if err := row.Scan(&positionsCount); err != nil {
		pr.log.Error("failed to scan positions count", zap.Error(err))

		return -1, fmt.Errorf("failed to scan positions count: %w", err)
	}

	return positionsCount, nil
}

// GetPositions returns a slice of positions for the given domain.
func (pr *PositionRepo) GetPositions(ctx context.Context, domain, orderBy string, limit, offset int) ([]*Position, error) {
	// Set default order in case if empty is given
	if orderBy == "" {
		orderBy = "volume"
	}

	rows, err := pr.conn.QueryContext(ctx,
		fmt.Sprintf(getPositionsQuery, orderBy),
		domain,
		limit,
		offset,
	)
	if err != nil {
		pr.log.Error("failed to execute query", zap.Error(err))

		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	positions := make([]*Position, 0, limit)

	for rows.Next() {
		p := &Position{}
		if err := rows.Scan(&p.Keyword,
			&p.Position,
			&p.URL,
			&p.Volume,
			&p.Results,
			&p.CPC,
			&p.Updated); err != nil {
			pr.log.Error("failed to scan position", zap.Error(err))

			return nil, fmt.Errorf("failed to scan position: %w", err)
		}
		positions = append(positions, p)
	}

	if err := rows.Err(); err != nil {
		pr.log.Error("failed to execute query", zap.Error(err))

		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return positions, nil
}
