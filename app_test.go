package qrcodeapi

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestServer(ctx context.Context, r router) *httptest.Server {
	e := newEcho()
	r.Route(e, "")

	ts := httptest.NewServer(e)
	go func() {
		<-ctx.Done()
		ts.Close()
	}()

	return ts
}

func TestParseInt(t *testing.T) {
	type args struct {
		value        string
		defaultValue int
		minValue     int
		maxValue     int
	}
	tests := [...]struct {
		name string
		args args
		want int
	}{
		{"default value", args{"", 10, 1, 100}, 10},
		{"cut max", args{"200", 10, 1, 100}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, parseIntDef(tt.args.value, tt.args.defaultValue, tt.args.minValue, tt.args.maxValue))
		})
	}
}
