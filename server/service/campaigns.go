package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	"github.com/igm/sockjs-go/sockjs"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/websocket"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) NewDistributedQueryCampaign(ctx context.Context, queryString string, hosts []uint, labels []uint) (*kolide.DistributedQueryCampaign, error) {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}

	query, err := svc.ds.NewQuery(&kolide.Query{
		Name:     fmt.Sprintf("distributed_%s_%d", vc.Username(), time.Now().Unix()),
		Query:    queryString,
		Saved:    false,
		AuthorID: vc.UserID(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "new query")
	}

	campaign, err := svc.ds.NewDistributedQueryCampaign(&kolide.DistributedQueryCampaign{
		QueryID: query.ID,
		Status:  kolide.QueryWaiting,
		UserID:  vc.UserID(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "new campaign")
	}

	// Add host targets
	for _, hid := range hosts {
		_, err = svc.ds.NewDistributedQueryCampaignTarget(&kolide.DistributedQueryCampaignTarget{
			Type: kolide.TargetHost,
			DistributedQueryCampaignID: campaign.ID,
			TargetID:                   hid,
		})
		if err != nil {
			return nil, errors.Wrap(err, "adding host target")
		}
	}

	// Add label targets
	for _, lid := range labels {
		_, err = svc.ds.NewDistributedQueryCampaignTarget(&kolide.DistributedQueryCampaignTarget{
			Type: kolide.TargetLabel,
			DistributedQueryCampaignID: campaign.ID,
			TargetID:                   lid,
		})
		if err != nil {
			return nil, errors.Wrap(err, "adding label target")
		}
	}

	return campaign, nil
}

type targetTotals struct {
	Total           uint `json:"count"`
	Online          uint `json:"online"`
	Offline         uint `json:"offline"`
	MissingInAction uint `json:"missing_in_action"`
}

const (
	campaignStatusPending  = "pending"
	campaignStatusFinished = "finished"
)

type campaignStatus struct {
	ExpectedResults uint   `json:"expected_results"`
	ActualResults   uint   `json:"actual_results"`
	Status          string `json:"status"`
}

func (svc service) StreamCampaignResults(ctx context.Context, conn *websocket.Conn, campaignID uint) {
	// Find the campaign and ensure it is active
	campaign, err := svc.ds.DistributedQueryCampaign(campaignID)
	if err != nil {
		conn.WriteJSONError(fmt.Sprintf("cannot find campaign for ID %d", campaignID))
		return
	}

	if campaign.Status != kolide.QueryWaiting {
		conn.WriteJSONError(fmt.Sprintf("campaign %d not running", campaignID))
		return
	}

	// Setting status to running will cause the query to be returned to the
	// targets when they check in for their queries
	campaign.Status = kolide.QueryRunning
	if err := svc.ds.SaveDistributedQueryCampaign(campaign); err != nil {
		conn.WriteJSONError("error saving campaign state")
		return
	}

	// Setting the status to completed stops the query from being sent to
	// targets. If this fails, there is a background job that will clean up
	// this campaign.
	defer func() {
		campaign.Status = kolide.QueryComplete
		svc.ds.SaveDistributedQueryCampaign(campaign)
	}()

	// Open the channel from which we will receive incoming query results
	// (probably from the redis pubsub implementation)
	readChan, err := svc.resultStore.ReadChannel(context.Background(), *campaign)
	if err != nil {
		conn.WriteJSONError(fmt.Sprintf("cannot open read channel for campaign %d ", campaignID))
		return
	}

	status := campaignStatus{
		Status: campaignStatusPending,
	}

	lastStatus := status.Status

	// to improve performance of the frontend rendering the results table, we
	// add the "host_hostname" field to every row.
	mapHostnameRows := func(hostname string, rows []map[string]string) {
		for _, row := range rows {
			row["host_hostname"] = hostname
		}
	}

	// Loop, pushing updates to results and expected totals
	for {
		// Update the expected hosts total (Should happen before
		// any results are written, to avoid the frontend showing "x of
		// 0 Hosts Returning y Records")
		hostIDs, labelIDs, err := svc.ds.DistributedQueryCampaignTargetIDs(campaign.ID)
		if err != nil {
			if err = conn.WriteJSONError("error retrieving campaign targets"); err != nil {
				return
			}
		}

		metrics, err := svc.CountHostsInTargets(context.Background(), hostIDs, labelIDs)
		if err != nil {
			if err = conn.WriteJSONError("error retrieving target counts"); err != nil {
				return
			}
		}

		totals := targetTotals{
			Total:           metrics.TotalHosts,
			Online:          metrics.OnlineHosts,
			Offline:         metrics.OfflineHosts,
			MissingInAction: metrics.MissingInActionHosts,
		}
		if err = conn.WriteJSONMessage("totals", totals); err != nil {
			return
		}

		status.ExpectedResults = totals.Online
		if status.ActualResults >= status.ExpectedResults {
			status.Status = campaignStatusFinished
		}
		// only write status message if status has changed
		if lastStatus != status.Status {
			lastStatus = status.Status
			if err = conn.WriteJSONMessage("status", status); err != nil {
				return
			}
		}
		select {
		case res := <-readChan:
			// Receive a result and push it over the websocket
			switch res := res.(type) {
			case kolide.DistributedQueryResult:
				mapHostnameRows(res.Host.HostName, res.Rows)
				err = conn.WriteJSONMessage("result", res)
				if err != nil {
					svc.logger.Log("msg", "error writing to channel", "err", err)
				}
				status.ActualResults++
			}

		case <-time.After(1 * time.Second):
			// Back to top of loop to update host totals
		}
	}

}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeCreateDistributedQueryCampaignRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createDistributedQueryCampaignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

// Create Distributed Query Campaign
type createDistributedQueryCampaignRequest struct {
	Query    string `json:"query"`
	Selected struct {
		Labels []uint `json:"labels"`
		Hosts  []uint `json:"hosts"`
	} `json:"selected"`
}

type createDistributedQueryCampaignResponse struct {
	Campaign *kolide.DistributedQueryCampaign `json:"campaign,omitempty"`
	Err      error                            `json:"error,omitempty"`
}

func (r createDistributedQueryCampaignResponse) error() error { return r.Err }

func makeCreateDistributedQueryCampaignEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createDistributedQueryCampaignRequest)
		campaign, err := svc.NewDistributedQueryCampaign(ctx, req.Query, req.Selected.Hosts, req.Selected.Labels)
		if err != nil {
			return createQueryResponse{Err: err}, nil
		}
		return createDistributedQueryCampaignResponse{campaign, nil}, nil
	}
}

// Stream Distributed Query Campaign Results and Metadata
func makeStreamDistributedQueryCampaignResultsHandler(svc kolide.Service, jwtKey string, logger kitlog.Logger) http.Handler {
	opt := sockjs.DefaultOptions
	opt.Websocket = true
	return sockjs.NewHandler("/api/v1/kolide/results", opt, func(session sockjs.Session) {
		defer session.Close(0, "none")

		conn := &websocket.Conn{Session: session}

		// Receive the auth bearer token
		token, err := conn.ReadAuthToken()
		if err != nil {
			logger.Log("err", err, "msg", "failed to read auth token")
			return
		}

		// Authenticate with the token
		vc, err := authViewer(context.Background(), jwtKey, token, svc)
		if err != nil || !vc.CanPerformActions() {
			logger.Log("err", err, "msg", "unauthorized viewer")
			conn.WriteJSONError("unauthorized")
			return
		}

		ctx := viewer.NewContext(context.Background(), *vc)

		msg, err := conn.ReadJSONMessage()
		if err != nil {
			logger.Log("err", err, "msg", "reading select_campaign JSON")
			conn.WriteJSONError("error reading select_campaign")
			return
		}
		if msg.Type != "select_campaign" {
			logger.Log("err", "unexpected msg type, expected select_campaign", "msg-type", msg.Type)
			conn.WriteJSONError("expected select_campaign")
			return
		}

		var info struct {
			CampaignID uint `json:"campaign_id"`
		}
		err = json.Unmarshal(*(msg.Data.(*json.RawMessage)), &info)
		if err != nil {
			logger.Log("err", err, "msg", "unmarshaling select_campaign data")
			conn.WriteJSONError("error unmarshaling select_campaign data")
			return
		}
		if info.CampaignID == 0 {
			logger.Log("err", "campaign ID not set")
			conn.WriteJSONError("0 is not a valid campaign ID")
			return
		}

		svc.StreamCampaignResults(ctx, conn, info.CampaignID)

	})
}
