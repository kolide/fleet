package launcher

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/kolide/osquery-go/plugin/distributed"
	"github.com/kolide/osquery-go/plugin/logger"
	"github.com/pkg/errors"

	"github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/kolide"
)

// launcherWrapper wraps the TLS interface.
type launcherWrapper struct {
	tls kolide.OsqueryService
}

func (svc *launcherWrapper) RequestEnrollment(ctx context.Context, enrollSecret, hostIdentifier string) (string, bool, error) {
	nodeKey, err := svc.tls.EnrollAgent(ctx, enrollSecret, hostIdentifier)
	if err != nil {
		if authErr, ok := err.(nodeInvalidErr); ok {
			return "", authErr.NodeInvalid(), err
		}
		return "", false, err
	}
	return nodeKey, false, nil
}

func (svc *launcherWrapper) RequestConfig(ctx context.Context, nodeKey string) (string, bool, error) {
	newCtx, invalid, err := svc.authenticateHost(ctx, nodeKey)
	if err != nil {
		return "", invalid, err
	}

	config, err := svc.tls.GetClientConfig(newCtx)
	if err != nil {
		return "", false, errors.Wrap(err, "get config for launcher")
	}

	// Launcher manages plugins so remove them from configuration if they exist.
	for _, optionName := range []string{"distributed_plugin", "logger_plugin"} {
		if _, ok := config.Options[optionName]; ok {
			delete(config.Options, optionName)
		}
	}

	buf := new(bytes.Buffer)
	if err = json.NewEncoder(buf).Encode(config); err != nil {
		return "", false, errors.Wrap(err, "encoding config for launcher")
	}

	return buf.String(), false, nil
}

func (svc *launcherWrapper) RequestQueries(ctx context.Context, nodeKey string) (*distributed.GetQueriesResult, bool, error) {
	newCtx, invalid, err := svc.authenticateHost(ctx, nodeKey)
	if err != nil {
		return nil, invalid, err
	}

	queryMap, accelerate, err := svc.tls.GetDistributedQueries(newCtx)
	if err != nil {
		return nil, false, errors.Wrap(err, "get queries for launcher")
	}

	result := &distributed.GetQueriesResult{
		Queries:           queryMap,
		AccelerateSeconds: int(accelerate),
	}

	return result, false, err
}

// TODO need to implement this method
func (svc *launcherWrapper) PublishLogs(ctx context.Context, nodeKey string, logType logger.LogType, logs []string) (string, string, bool, error) {
	newCtx, invalid, err := svc.authenticateHost(ctx, nodeKey)
	if err != nil {
		return "", "", invalid, err
	}
	_ = newCtx

	switch logType {
	case logger.LogTypeStatus:
	case logger.LogTypeSnapshot:
	default:
		panic("chat with @zwass about what to do with the rest")
	}
	panic("todo")
}

func (svc *launcherWrapper) PublishResults(ctx context.Context, nodeKey string, results []distributed.Result) (string, string, bool, error) {
	newCtx, invalid, err := svc.authenticateHost(ctx, nodeKey)
	if err != nil {
		return "", "", invalid, err
	}

	osqueryResults := make(kolide.OsqueryDistributedQueryResults, len(results))
	statuses := make(map[string]string, len(results))

	for _, result := range results {
		statuses[result.QueryName] = strconv.Itoa(result.Status)
		osqueryResults[result.QueryName] = result.Rows
	}

	if err := svc.tls.SubmitDistributedQueryResults(newCtx, osqueryResults, statuses); err != nil {
		return "", "", false, errors.Wrap(err, "submit launcher results")
	}

	return "", "", false, nil
}

func (svc *launcherWrapper) CheckHealth(ctx context.Context) (int32, error) {
	// TODO we can pass the healthcheckers as optional parameters during init and implement this.
	return 0, nil
}

// AuthenticateHost verifies the host node key using the TLS API and returns back a
// context which includes the host as a context value.
// In the kolide.OsqueryService authentication is done via endpoint middleware, but all launcher endpoints require
// an explicit return for NodeInvalid, so we check in this helper method instead.
// TODO(@groob) we can probably expose the kolide endpoint middleware to do the validation in the endpoint layer here
// as well... at least document and come back to this in another PR.
func (svc *launcherWrapper) authenticateHost(ctx context.Context, nodeKey string) (context.Context, bool, error) {
	node, err := svc.tls.AuthenticateHost(ctx, nodeKey)
	if err != nil {
		if authErr, ok := err.(nodeInvalidErr); ok {
			return ctx, authErr.NodeInvalid(), err
		}
		return ctx, false, err
	}

	ctx = host.NewContext(ctx, *node)
	return ctx, false, nil
}

type nodeInvalidErr interface {
	error
	NodeInvalid() bool
}
