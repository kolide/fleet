package kolide

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type QueryStore interface {
	// Query methods
	NewQuery(query *Query) (*Query, error)
	SaveQuery(query *Query) error
	DeleteQuery(query *Query) error
	Query(id uint) (*Query, error)
	ListQueries(opt ListOptions) ([]*Query, error)

	// NewDistributedQueryCampaign creates a new distributed query campaign
	NewDistributedQueryCampaign(camp *DistributedQueryCampaign) (*DistributedQueryCampaign, error)
	// SaveDistributedQueryCampaign updates an existing distributed query
	// campaign
	SaveDistributedQueryCampaign(camp *DistributedQueryCampaign) error
	// NewDistributedQueryCampaignTarget adds a new target to an existing
	// distributed query campaign
	NewDistributedQueryCampaignTarget(target *DistributedQueryCampaignTarget) (*DistributedQueryCampaignTarget, error)
	// NewDistributedQueryCampaignExecution records a new execution for a
	// distributed query campaign
	NewDistributedQueryExecution(exec *DistributedQueryExecution) (*DistributedQueryExecution, error)
}

type QueryService interface {
	ListQueries(ctx context.Context, opt ListOptions) ([]*Query, error)
	GetQuery(ctx context.Context, id uint) (*Query, error)
	NewQuery(ctx context.Context, p QueryPayload) (*Query, error)
	ModifyQuery(ctx context.Context, id uint, p QueryPayload) (*Query, error)
	DeleteQuery(ctx context.Context, id uint) error
	NewDistributedQueryCampaign(ctx context.Context, queryString string, hosts []uint, labels []uint) (*DistributedQueryCampaign, error)

	// StreamCampaignResults is a handler function that upgrades a request
	// to a websocket connection and begins streaming distributed query
	// results. Note that it uses a type signature inconsistent with other
	// service methods due to go-kit's lack of support of websockets/pubsub
	StreamCampaignResults(jwtKey string) func(w http.ResponseWriter, r *http.Request)
}

type QueryPayload struct {
	Name         *string
	Description  *string
	Query        *string
	Interval     *uint
	Snapshot     *bool
	Differential *bool
	Platform     *string
	Version      *string
}

type Query struct {
	UpdateCreateTimestamps
	DeleteFields
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Query        string `json:"query"`
	Interval     uint   `json:"interval"`
	Snapshot     bool   `json:"snapshot"`
	Differential bool   `json:"differential"`
	Platform     string `json:"platform"`
	Version      string `json:"version"`
}

type DistributedQueryStatus int

const (
	QueryRunning  DistributedQueryStatus = iota
	QueryComplete DistributedQueryStatus = iota
	QueryError    DistributedQueryStatus = iota
)

type DistributedQueryCampaign struct {
	UpdateCreateTimestamps
	DeleteFields
	ID      uint                   `json:"id"`
	QueryID uint                   `json:"query_id" db:"query_id"`
	Status  DistributedQueryStatus `json:"status"`
	UserID  uint                   `json:"user_id"`
}

type DistributedQueryCampaignTarget struct {
	ID                         uint
	Type                       TargetType
	DistributedQueryCampaignID uint `db:"distributed_query_campaign_id"`
	TargetID                   uint `db:"target_id"`
}

type DistributedQueryExecutionStatus int

const (
	ExecutionWaiting DistributedQueryExecutionStatus = iota
	ExecutionRequested
	ExecutionSucceeded
	ExecutionFailed
)

type DistributedQueryResult struct {
	DistributedQueryCampaignID uint                `json:"distributed_query_execution_id"`
	Host                       Host                `json:"host"`
	Rows                       []map[string]string `json:"rows"`
}

type DistributedQueryExecution struct {
	ID                         uint
	HostID                     uint `db:"host_id"`
	DistributedQueryCampaignID uint `db:"distributed_query_campaign_id"`
	Status                     DistributedQueryExecutionStatus
	Error                      string
	ExecutionDuration          time.Duration `db:"execution_duration"`
}

type Option struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Key       string
	Value     string
	Platform  string
}

type DecoratorType int

const (
	DecoratorLoad DecoratorType = iota
	DecoratorAlways
	DecoratorInterval
)

type Decorator struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Type      DecoratorType
	Interval  int
	Query     string
}
