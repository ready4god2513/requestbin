package store

import (
	"context"
	"encoding/json"

	"github.com/ready4god2513/requestbin/model"
)

func (s *Store) CreateRequest(ctx context.Context, req *model.Request) (*model.Request, error) {
	headersJSON, err := json.Marshal(req.Headers)
	if err != nil {
		return nil, err
	}
	queryParamsJSON, err := json.Marshal(req.QueryParams)
	if err != nil {
		return nil, err
	}

	var (
		created        model.Request
		headersRaw     []byte
		queryParamsRaw []byte
	)

	err = s.pool.QueryRow(ctx, `
		INSERT INTO requests
			(bin_id, method, path, headers, query_params, body, remote_addr, content_type, content_length)
		VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6, $7, $8, $9)
		RETURNING
			id, bin_id, method, path, headers, query_params,
			COALESCE(body, ''), COALESCE(remote_addr, ''),
			COALESCE(content_type, ''), COALESCE(content_length, 0), created_at
	`,
		req.BinID, req.Method, req.Path,
		string(headersJSON), string(queryParamsJSON),
		req.Body, req.RemoteAddr, req.ContentType, req.ContentLength,
	).Scan(
		&created.ID, &created.BinID, &created.Method, &created.Path,
		&headersRaw, &queryParamsRaw,
		&created.Body, &created.RemoteAddr, &created.ContentType, &created.ContentLength,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(headersRaw, &created.Headers)
	_ = json.Unmarshal(queryParamsRaw, &created.QueryParams)

	return &created, nil
}

func (s *Store) ListRequests(ctx context.Context, binID string) ([]*model.Request, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT
			id, bin_id, method, path, headers, query_params,
			COALESCE(body, ''), COALESCE(remote_addr, ''),
			COALESCE(content_type, ''), COALESCE(content_length, 0), created_at
		FROM requests
		WHERE bin_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, binID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*model.Request, 0)
	for rows.Next() {
		var r model.Request
		var headersRaw, queryParamsRaw []byte
		if err := rows.Scan(
			&r.ID, &r.BinID, &r.Method, &r.Path,
			&headersRaw, &queryParamsRaw,
			&r.Body, &r.RemoteAddr, &r.ContentType, &r.ContentLength,
			&r.CreatedAt,
		); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(headersRaw, &r.Headers)
		_ = json.Unmarshal(queryParamsRaw, &r.QueryParams)
		result = append(result, &r)
	}

	return result, rows.Err()
}

func (s *Store) ClearRequests(ctx context.Context, binID string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM requests WHERE bin_id = $1`, binID)
	return err
}

func (s *Store) DeleteRequest(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM requests WHERE id = $1`, id)
	return err
}
