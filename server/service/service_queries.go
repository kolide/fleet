package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kolide/kolide-ose/server/contexts/viewer"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) ListQueries(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Query, error) {
	return svc.ds.ListQueries(opt)
}

func (svc service) GetQuery(ctx context.Context, id uint) (*kolide.Query, error) {
	return svc.ds.Query(id)
}

func (svc service) NewQuery(ctx context.Context, p kolide.QueryPayload) (*kolide.Query, error) {
	query := &kolide.Query{}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	if p.Interval != nil {
		query.Interval = *p.Interval
	}

	if p.Snapshot != nil {
		query.Snapshot = *p.Snapshot
	}

	if p.Differential != nil {
		query.Differential = *p.Differential
	}

	if p.Platform != nil {
		query.Platform = *p.Platform
	}

	if p.Version != nil {
		query.Version = *p.Version
	}

	query, err := svc.ds.NewQuery(query)
	if err != nil {
		return nil, err
	}
	return query, nil
}

func (svc service) ModifyQuery(ctx context.Context, id uint, p kolide.QueryPayload) (*kolide.Query, error) {
	query, err := svc.ds.Query(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	if p.Interval != nil {
		query.Interval = *p.Interval
	}

	if p.Snapshot != nil {
		query.Snapshot = *p.Snapshot
	}

	if p.Differential != nil {
		query.Differential = *p.Differential
	}

	if p.Platform != nil {
		query.Platform = *p.Platform
	}

	if p.Version != nil {
		query.Version = *p.Version
	}

	err = svc.ds.SaveQuery(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (svc service) DeleteQuery(ctx context.Context, id uint) error {
	query, err := svc.ds.Query(id)
	if err != nil {
		return err
	}

	err = svc.ds.DeleteQuery(query)
	if err != nil {
		return err
	}

	return nil
}

func (svc service) NewDistributedQueryCampaign(ctx context.Context, queryString string, hosts []uint, labels []uint) (*kolide.DistributedQueryCampaign, error) {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errNoContext
	}

	query, err := svc.NewQuery(ctx, kolide.QueryPayload{
		Name:  &queryString,
		Query: &queryString,
	})
	if err != nil {
		return nil, err
	}

	campaign, err := svc.ds.NewDistributedQueryCampaign(&kolide.DistributedQueryCampaign{
		QueryID: query.ID,
		Status:  kolide.QueryRunning,
		UserID:  vc.UserID(),
	})
	if err != nil {
		return nil, err
	}

	// Add host targets
	for _, hid := range hosts {
		_, err = svc.ds.NewDistributedQueryCampaignTarget(&kolide.DistributedQueryCampaignTarget{
			Type: kolide.TargetHost,
			DistributedQueryCampaignID: campaign.ID,
			TargetID:                   hid,
		})
		if err != nil {
			return nil, err
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
			return nil, err
		}
	}

	return campaign, nil
}

func (svc service) StreamCampaignResults(jwtKey string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer conn.Close()

		// Receive the auth bearer token
		_, token, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("here1")
		// Authenticate with the token
		vc, err := authViewer(context.Background(), jwtKey, string(token), svc)
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("unauthorized"))
			return
		}
		if !vc.CanPerformActions() {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("unauthorized"))
			return
		}

		fmt.Println("here2")
		campaignID, err := idFromRequest(r, "id")
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("invalid campaign ID"))
			return
		}

		fmt.Println("here3")
		readChan, err := svc.resultStore.ReadChannel(context.Background(), kolide.DistributedQueryCampaign{ID: campaignID})
		if err != nil {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("cannot open read channel for campaign %d ", campaignID)))
			return
		}

		for res := range readChan {
			switch res := res.(type) {
			case kolide.DistributedQueryResult:
				err = conn.WriteJSON(res)
				if err != nil {
					fmt.Println("error writing to channel")
				}
			}
		}

		for {
			if err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%d", time.Now()))); err != nil {
				fmt.Println(err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}
