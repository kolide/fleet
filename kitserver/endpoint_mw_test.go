package kitserver

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/datastore"
)

func TestViewerContext(t *testing.T) {
	ctx := context.Background()
	req := struct{}{}
	ds, _ := datastore.New("mock", "")
	createTestUsers(t, ds)
	admin1, _ := ds.User("admin1")

	e := endpoint.Nop // a test endpoint
	var endpointTests = []struct {
		endpoint endpoint.Endpoint
		vc       *ViewerContext
		err      error
		autherr  bool
	}{
		{
			endpoint: mustBeAdmin(e),
			err:      errNoContext,
		},
		{
			endpoint: mustBeAdmin(e),
			vc:       &ViewerContext{user: admin1},
		},
	}

	for _, tt := range endpointTests {
		if tt.vc != nil {
			ctx = context.WithValue(ctx, "viewerContext", tt.vc)
		}
		if _, err := tt.endpoint(ctx, req); err != tt.err {
			t.Fatal(err)
		}

	}

}
