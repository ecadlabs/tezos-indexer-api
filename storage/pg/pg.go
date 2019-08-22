package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/ecadlabs/tezos-indexer-api/storage"
	"github.com/jackc/pgx/v4"
)

const (
	defaultLimit = 1000
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type PostgresStorage struct {
	DB Queryer
}

func (p *PostgresStorage) GetBalanceUpdate(ctx context.Context, address string, start, end time.Time, limit int) ([]*storage.BalanceUpdate, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	query := `
		SELECT
		    level,
			timestamp,
			diff::numeric,
			SUM(diff::numeric) OVER (ORDER BY level) AS value
		FROM
			balance 
			JOIN block ON balance.block_hash = block.hash 
		WHERE
			contract_address = $1
		`

	arg := []interface{}{address}
	idx := 2

	if !end.IsZero() {
		query += fmt.Sprintf(" AND timestamp < $%d", idx)
		arg = append(arg, end)
		idx++
	}

	query += " ORDER BY level"

	if !start.IsZero() {
		// WHERE condition affects the window frame and the computation of the integral balance value so use a subquery
		query = fmt.Sprintf("SELECT * FROM (%s) AS bal WHERE timestamp >= $%d", query, idx)
		arg = append(arg, start)
		idx++
	}

	query += fmt.Sprintf(" LIMIT $%d", idx)
	arg = append(arg, limit)
	idx++

	rows, err := p.DB.Query(ctx, query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]*storage.BalanceUpdate, 0, limit)

	for rows.Next() {
		var v storage.BalanceUpdate
		err = rows.Scan(
			&v.BlockLevel,
			&v.BlockTimestamp,
			&v.Diff,
			&v.Value)

		if err != nil {
			return nil, err
		}

		res = append(res, &v)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return res, nil
}
