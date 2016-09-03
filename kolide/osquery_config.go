package kolide

import (
	"time"

	"golang.org/x/net/context"
)

type OsqueryConfigStore interface {
	// Query methods
	NewQuery(query *Query) error
	SaveQuery(query *Query) error
	DeleteQuery(query *Query) error
	Query(id uint) (*Query, error)
	Queries() ([]*Query, error)

	// Label methods
	NewLabel(label *Label) error

	// Pack methods
	NewPack(pack *Pack) error
	SavePack(pack *Pack) error
	DeletePack(pack *Pack) error
	Pack(id uint) (*Pack, error)
	Packs() ([]*Pack, error)

	// Modifying the queries in packs
	AddQueryToPack(query *Query, pack *Pack) error
	GetQueriesInPack(pack *Pack) ([]*Query, error)
	RemoveQueryFromPack(query *Query, pack *Pack) error
}

type OsqueryConfigService interface {
	GetAllQueries(ctx context.Context) ([]*Query, error)
	GetQuery(ctx context.Context, id uint) (*Query, error)
	NewQuery(ctx context.Context, p QueryPayload) (*Query, error)
	ModifyQuery(ctx context.Context, id uint, p QueryPayload) (*Query, error)
	DeleteQuery(ctx context.Context, id uint) error

	GetAllPacks(ctx context.Context) ([]*Pack, error)
	GetPack(ctx context.Context, id uint) (*Pack, error)
	NewPack(ctx context.Context, p PackPayload) (*Pack, error)
	ModifyPack(ctx context.Context, id uint, p PackPayload) (*Pack, error)
	DeletePack(ctx context.Context, id uint) error

	AddQueryToPack(ctx context.Context, qid, pid uint) error
	RemoveQueryFromPack(ctx context.Context, qid, pid uint) error
}

type QueryPayload struct {
	Name         *string
	Query        *string
	Interval     *uint
	Snapshot     *bool
	Differential *bool
	Platform     *string
	Version      *string
}

type PackPayload struct {
	Name     *string
	Platform *string
}

type Query struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Name         string `gorm:"not null;unique_index:idx_query_unique_name"`
	Query        string `gorm:"not null"`
	Interval     uint
	Snapshot     bool
	Differential bool
	Platform     string
	Version      string
}

type Label struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique_index:idx_label_unique_name"`
	QueryID   uint
}

type LabelQueryExecution struct {
	ID        uint `gorm:"primary_key"`
	UpdatedAt time.Time
	Matches   bool
	LabelID   uint // Note we manually specify a unique index on these
	HostID    uint // fields in gormDB.Migrate
}

type TargetType int

const (
	TargetLabel TargetType = iota
	TargetHost  TargetType = iota
)

type Target struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      TargetType
	TargetID  uint
	QueryID   uint
}

type Pack struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique_index:idx_pack_unique_name"`
	Platform  string
}

type PackQuery struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	PackID    uint
	QueryID   uint
}

type PackTarget struct {
	ID       uint `gorm:"primary_key"`
	PackID   uint
	TargetID uint
}

type DistributedQueryStatus int

const (
	QueryRunning  DistributedQueryStatus = iota
	QueryComplete DistributedQueryStatus = iota
	QueryError    DistributedQueryStatus = iota
)

type DistributedQueryCampaign struct {
	ID          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	QueryID     uint
	MaxDuration time.Duration
	Status      DistributedQueryStatus
	UserID      uint
}

type DistributedQueryCampaignTarget struct {
	ID                         uint `gorm:"primary_key"`
	DistributedQueryCampaignID uint
	TargetID                   uint
}

type DistributedQueryExecutionStatus int

const (
	ExecutionWaiting   DistributedQueryExecutionStatus = iota
	ExecutionRequested DistributedQueryExecutionStatus = iota
	ExecutionSucceeded DistributedQueryExecutionStatus = iota
	ExecutionFailed    DistributedQueryExecutionStatus = iota
)

type DistributedQueryExecution struct {
	HostID             uint
	DistributedQueryID uint
	Status             DistributedQueryExecutionStatus
	Error              string `gorm:"size:1024"`
	ExecutionDuration  time.Duration
}

type Option struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string `gorm:"not null;unique_index:idx_option_unique_key"`
	Value     string `gorm:"not null"`
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad     DecoratorType = iota
	DecoratorAlways   DecoratorType = iota
	DecoratorInterval DecoratorType = iota
)

type Decorator struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType `gorm:"not null"`
	Interval  int
	Query     string
}
