package testutils

import (
	"context"
	"net/http/httptest"

	"qrcodeapi/pkg/helper/echox"
)

func NewTestServer(ctx context.Context, r echox.Router) *httptest.Server {
	e := echox.New()
	r.Route(e, "")

	ts := httptest.NewServer(e)
	go func() {
		<-ctx.Done()
		ts.Close()
	}()

	return ts
}
