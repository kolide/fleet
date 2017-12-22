package kolide

import (
	"context"
)

type ScheduledQueryStore interface {
	ListScheduledQueriesInPack(id uint, opts ListOptions) ([]*ScheduledQuery, error)
}

type ScheduledQueryService interface {
	GetScheduledQueriesInPack(ctx context.Context, id uint, opts ListOptions) (queries []*ScheduledQuery, err error)
}

type ScheduledQuery struct {
	UpdateCreateTimestamps
	DeleteFields
	ID       uint    `json:"id"`
	PackID   uint    `json:"pack_id" db:"pack_id"`
	QueryID  uint    `json:"query_id" db:"query_id"`
	Query    string  `json:"query"` // populated via a join on queries
	Name     string  `json:"name"`  // populated via a join on queries
	Interval uint    `json:"interval"`
	Snapshot *bool   `json:"snapshot"`
	Removed  *bool   `json:"removed"`
	Platform *string `json:"platform"`
	Version  *string `json:"version"`
	Shard    *uint   `json:"shard"`
}

type ScheduledQueryPayload struct {
	PackID   *uint   `json:"pack_id"`
	QueryID  *uint   `json:"query_id"`
	Interval *uint   `json:"interval"`
	Snapshot *bool   `json:"snapshot"`
	Removed  *bool   `json:"removed"`
	Platform *string `json:"platform"`
	Version  *string `json:"version"`
	Shard    *uint   `json:"shard"`
}
