package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/jackc/pgx/v5"
	"github.com/ready4god2513/requestbin/model"
)

func (s *Store) CreateBin(ctx context.Context) (*model.Bin, error) {
	id, err := generateID()
	if err != nil {
		return nil, err
	}

	var bin model.Bin
	err = s.pool.QueryRow(ctx,
		`INSERT INTO bins (id) VALUES ($1) RETURNING id, created_at`,
		id,
	).Scan(&bin.ID, &bin.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &bin, nil
}

func (s *Store) GetBin(ctx context.Context, id string) (*model.Bin, error) {
	var bin model.Bin
	err := s.pool.QueryRow(ctx, `
		SELECT b.id, b.created_at, COUNT(r.id)
		FROM bins b
		LEFT JOIN requests r ON r.bin_id = b.id
		WHERE b.id = $1
		GROUP BY b.id, b.created_at
	`, id).Scan(&bin.ID, &bin.CreatedAt, &bin.RequestCount)
	if err == pgx.ErrNoRows {
		return nil, pgx.ErrNoRows
	}
	return &bin, err
}

func (s *Store) DeleteBin(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM bins WHERE id = $1`, id)
	return err
}

func generateID() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
