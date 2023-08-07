package apiserver

import (
	"context"
	"testing"
	"time"

	"qrcodeapi/config"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/request"
)

func TestApp(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := New()
	go service.Serve(ctx)
	time.Sleep(time.Second)

	addr := config.BindAddr()

	resp, err := request.Get("http://%s/api/v1/qrcode?content=HELLO", addr).Do(ctx)
	require.NoError(t, err)
	require.NoErrorf(t, resp.Success(), "failed with status %d", resp.StatusCode)

	require.Equal(t, "image/png", resp.Header.Get(echo.HeaderContentType))
}
