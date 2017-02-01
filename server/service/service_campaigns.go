package service

import (
	"fmt"
	"time"

	"github.com/kolide/kolide/server/contexts/viewer"
	"github.com/kolide/kolide/server/kolide"
	"github.com/kolide/kolide/server/websocket"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

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
			// Update the expected hosts total
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

		}
	}

}
