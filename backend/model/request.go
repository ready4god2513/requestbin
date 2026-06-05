package model

import "time"

type Request struct {
	ID            string            `json:"id"`
	BinID         string            `json:"bin_id"`
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Headers       map[string]string `json:"headers"`
	QueryParams   map[string]string `json:"query_params"`
	Body          string            `json:"body"`
	RemoteAddr    string            `json:"remote_addr"`
	ContentType   string            `json:"content_type"`
	ContentLength int64             `json:"content_length"`
	CreatedAt     time.Time         `json:"created_at"`
}
