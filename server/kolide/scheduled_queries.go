package kolide

import (
	"golang.org/x/net/context"
)

type ScheduledQueryStore interface {
	NewScheduledQuery(sq *PackQuery) (*PackQuery, error)
	SaveScheduledQuery(sq *PackQuery) (*PackQuery, error)
	DeleteScheduledQuery(id uint) error
	ScheduledQuery(id uint) (*PackQuery, error)
	ListScheduledQueriesInPack(id uint, opts ListOptions) ([]*PackQuery, error)
}

type ScheduledQueryService interface {
	GetScheduledQuery(ctx context.Context, id uint) (*PackQuery, error)
	GetScheduledQueriesInPack(ctx context.Context, id uint, opts ListOptions) ([]*PackQuery, error)
	ScheduleQuery(ctx context.Context, sq *PackQuery) (*PackQuery, error)
	DeleteScheduledQuery(ctx context.Context, id uint) error
	ModifyScheduledQuery(ctx context.Context, sq *PackQuery) (*PackQuery, error)
}

type PackQuery struct {
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
