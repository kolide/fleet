package service

import (
	"encoding/json"
	"net/http"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type importConfigRequest struct {
	// Config contains a JSON osquery config supplied by the end user
	Config string `json:"config"`
	// ExternalPackConfigs contains a map of external Pack configs keyed by
	// Pack name, this includes external packs referenced by the globbing
	// feature.  Not in the case of globbed packs, we expect the user to
	// generate unique pack names since we don't know what they are, these
	// names must be included in the GlobPackNames field so that we can
	// validate that they've been accounted for.
	ExternalPackConfigs map[string]string `json:"external_pack_configs"`
	// GlobPackNames list of user generated names for external packs
	// referenced by the glob feature, the JSON for the globbed packs
	// is stored in ExternalPackConfigs keyed by the GlobPackName
	GlobPackNames []string `json:"glob_pack_names"`
}

func decodeImportConfigRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req importConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	// Unmarshal main config
	conf := kolide.ImportConfig{
		Packs:         make(kolide.PackNameMap),
		ExternalPacks: make(kolide.PackNameToPackDetails),
	}
	if err := json.Unmarshal([]byte(req.Config), &conf); err != nil {
		return nil, err
	}
	// Unmarshal external packs
	for packName, packConfig := range req.ExternalPackConfigs {
		var pack kolide.PackDetails
		if err := json.Unmarshal([]byte(packConfig), &pack); err != nil {
			return nil, err
		}
		conf.ExternalPacks[packName] = pack
	}
	conf.GlobPackNames = req.GlobPackNames
	return conf, nil
}
