package launcher

import (
	"bytes"
	newcontext "context"
	"encoding/json"

	pb "github.com/kolide/agent-api"
	"github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

var errNotImplmented = errors.New("not implemented")

// agentBinding implements ApiClient interface and maps gRPC domain functions to the application.
type agentBinding struct {
	service kolide.OsqueryService
}

func newAgentBinding(svc kolide.OsqueryService) pb.ApiServer {
	return &agentBinding{
		service: svc,
	}
}

type enrollmentError interface {
	NodeInvalid() bool
	Error() string
}

// Attempt to enroll a host with kolide/cloud
func (b *agentBinding) RequestEnrollment(ctx context.Context, req *pb.EnrollmentRequest) (*pb.EnrollmentResponse, error) {
	var resp pb.EnrollmentResponse
	nodeKey, err := b.service.EnrollAgent(newCtx(ctx), req.EnrollSecret, req.HostIdentifier)
	if err != nil {
		if errEnroll, ok := err.(enrollmentError); ok {
			resp.NodeInvalid = errEnroll.NodeInvalid()
			resp.ErrorCode = errEnroll.Error()
			return &resp, nil
		}
		return nil, err
	}
	resp.NodeKey = nodeKey
	return &resp, nil
}

// RequestConfig requests an updated configuration
func (b *agentBinding) RequestConfig(ctx context.Context, req *pb.AgentApiRequest) (*pb.ConfigResponse, error) {
	config, err := b.service.GetClientConfig(newCtx(ctx))
	if err != nil {
		return nil, err
	}
	var writer bytes.Buffer
	if err = json.NewEncoder(&writer).Encode(config); err != nil {
		return nil, err
	}
	return &pb.ConfigResponse{ConfigJsonBlob: writer.String()}, nil
}

// RequestQueries request/pull distributed queries
func (b *agentBinding) RequestQueries(ctx context.Context, _ *pb.AgentApiRequest) (*pb.QueryCollection, error) {
	queryMap, _, err := b.service.GetDistributedQueries(newCtx(ctx))
	if err != nil {
		return nil, err
	}
	var result pb.QueryCollection
	for id, query := range queryMap {
		result.Queries = append(result.Queries, &pb.QueryCollection_Query{Id: id, Query: query})
	}
	return &result, nil
}

type StatusLog struct {
	Severity string `json:"s"`
	Filename string `json:"f"`
	Line     string `json:"i"`
	Message  string `json:"m"`
}

// convert the json from grpc client to an object suitable
// for consumption by fleet
func toKolideLog(jsn string) (*kolide.OsqueryStatusLog, error) {
	var status StatusLog
	err := json.NewDecoder(bytes.NewBufferString(jsn)).Decode(&status)
	if err != nil {
		return nil, err
	}
	result := &kolide.OsqueryStatusLog{
		Severity: status.Severity,
		Filename: status.Filename,
		Line:     status.Line,
		Message:  status.Message,
	}
	return result, nil
}

// PublishLogs publish logs from osqueryd
func (b *agentBinding) PublishLogs(ctx context.Context, coll *pb.LogCollection) (*pb.AgentApiResponse, error) {
	if coll.LogType == pb.LogCollection_STATUS {
		var statuses []kolide.OsqueryStatusLog
		for _, record := range coll.Logs {

			status, err := toKolideLog(record.Data)
			if err != nil {
				return nil, errors.Wrap(err, "decoding status log")
			}
			statuses = append(statuses, *status)
		}

		if err := b.service.SubmitStatusLogs(newCtx(ctx), statuses); err != nil {
			return nil, errors.Wrap(err, "submitting status logs")
		}

	}
	return &pb.AgentApiResponse{}, nil
}

// PublishResults publish distributed query results
func (b *agentBinding) PublishResults(ctx context.Context, coll *pb.ResultCollection) (*pb.AgentApiResponse, error) {
	return &pb.AgentApiResponse{}, nil
}

// HotConfigure pushed configurations
func (b *agentBinding) HotConfigure(in *pb.AgentApiRequest, svr pb.Api_HotConfigureServer) error {
	return errNotImplmented
}

// HotlineBling this would be live query push to agent
func (b *agentBinding) HotlineBling(svr pb.Api_HotlineBlingServer) error {
	return errNotImplmented
}

// newCtx is used to map the old golang.com/net/context which we are forced to use
// because our generated gRPC code uses it, to the new stdlib context, which is used
// by the Fleet application.
func newCtx(ctx context.Context) newcontext.Context {
	if h, ok := ctx.Value(hostKey).(kolide.Host); ok {
		return host.NewContext(newcontext.Background(), h)
	}
	return newcontext.Background()
}
