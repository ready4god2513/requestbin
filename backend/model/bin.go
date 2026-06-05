package model

import "time"

type Bin struct {
	ID           string    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	RequestCount int64     `json:"request_count"`
}
