package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS bins (
			id         VARCHAR(12) PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS requests (
			id             UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
			bin_id         VARCHAR(12) NOT NULL REFERENCES bins(id) ON DELETE CASCADE,
			method         VARCHAR(10) NOT NULL,
			path           TEXT        NOT NULL DEFAULT '/',
			headers        JSONB       NOT NULL DEFAULT '{}',
			query_params   JSONB       NOT NULL DEFAULT '{}',
			body           TEXT,
			remote_addr    TEXT,
			content_type   TEXT,
			content_length BIGINT,
			created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_requests_bin_id     ON requests(bin_id);
		CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(created_at DESC);
	`)
	return err
}
