package testutils

import (
	"context"
	"net/http/httptest"

	"github.com/whitekid/echox"
)

func NewTestServer(ctx context.Context, r echox.Router) *httptest.Server {
	e := echox.New()
	e.Route(nil, r)

	ts := httptest.NewServer(e)
	go func() {
		<-ctx.Done()
		ts.Close()
	}()

	return ts
}
