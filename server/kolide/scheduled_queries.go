package kolide

import (
	"golang.org/x/net/context"
)

type ScheduledQueryStore interface {
	NewScheduledQuery(sq *ScheduledQuery) (*ScheduledQuery, error)
	SaveScheduledQuery(sq *ScheduledQuery) (*ScheduledQuery, error)
	DeleteScheduledQuery(id uint) error
	ScheduledQuery(id uint) (*ScheduledQuery, error)
	ListScheduledQueriesInPack(id uint, opts ListOptions) ([]*ScheduledQuery, error)
}

type ScheduledQueryService interface {
	GetScheduledQuery(ctx context.Context, id uint) (*ScheduledQuery, error)
	GetScheduledQueriesInPack(ctx context.Context, id uint, opts ListOptions) ([]*ScheduledQuery, error)
	ScheduleQuery(ctx context.Context, sq *ScheduledQuery) (*ScheduledQuery, error)
	DeleteScheduledQuery(ctx context.Context, id uint) error
	ModifyScheduledQuery(ctx context.Context, sq *ScheduledQuery) (*ScheduledQuery, error)
}

type ScheduledQuery struct {
	UpdateCreateTimestamps
	DeleteFields
	ID           uint
	PackID       uint    `json:"pack_id" db:"pack_id"`
	QueryID      uint    `json:"query_id" db:"query_id"`
	Interval     uint    `json:"interval"`
	Snapshot     *bool   `json:"snapshot"`
	Differential *bool   `json:"differential"`
	Platform     *string `json:"platform"`
	Version      *string `json:"version"`
	Shard        *uint   `json:"shard"`
}
