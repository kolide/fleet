package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Search Targrets
////////////////////////////////////////////////////////////////////////////////

type searchTargetsRequest struct {
	Query    string `json:"query"`
	Selected struct {
		Labels []uint `json:"labels"`
		Hosts  []uint `json:"hosts"`
	} `json:"selected"`
}

type targetsData struct {
	Hosts  []hostResponse `json:"hosts"`
	Labels []kolide.Label `json:"labels"`
}

type searchTargetsResponse struct {
	Targets              *targetsData `json:"targets,omitempty"`
	SelectedTargetsCount uint         `json:"selected_targets_count,omitempty"`
	Err                  error        `json:"error,omitempty"`
}

func (r searchTargetsResponse) error() error { return r.Err }

func makeSearchTargetsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(searchTargetsRequest)

		results, count, err := svc.SearchTargets(ctx, req.Query, req.Selected.Hosts, req.Selected.Labels)
		if err != nil {
			return searchTargetsResponse{Err: err}, nil
		}

		targets := &targetsData{
			Hosts:  []hostResponse{},
			Labels: []kolide.Label{},
		}

		for _, host := range results.Hosts {
			targets.Hosts = append(targets.Hosts, hostResponse{host, svc.HostStatus(ctx, host)})
		}

		for _, label := range results.Labels {
			targets.Labels = append(targets.Labels, label)
		}

		return searchTargetsResponse{
			Targets:              targets,
			SelectedTargetsCount: count,
		}, nil
	}
}
