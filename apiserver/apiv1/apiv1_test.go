package apiv1

import (
	"context"
	"image"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp/request"

	"qrcodeapi/pkg/ical"
	"qrcodeapi/pkg/qrcode"
	"qrcodeapi/pkg/testutils.go"
)

func TestText(t *testing.T) {
	type args struct {
		text string
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"default", args{"동해물과 백두산이"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ts := testutils.NewTestServer(ctx, NewAPIv1())

			req := request.Get("%s/api/v1/qrcode", ts.URL).Query("content", tt.args.text)

			resp, err := req.Do(ctx)
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Truef(t, resp.Success(), "failed with status %v", resp.StatusCode)

			img, _, err := image.Decode(resp.Body)
			require.NoError(t, err)

			got, err := qrcode.Decode(img)
			require.NoError(t, err)
			require.Equal(t, tt.args.text, got)
		})
	}
}

func TestAccept(t *testing.T) {
	type args struct {
		accept string
	}
	tests := [...]struct {
		name            string
		args            args
		wantStatus      int
		wantContentType string
		wantImage       string
	}{
		{"missing accept", args{""}, http.StatusOK, "image/png", "png"},
		{"curl/wget default", args{"*/*"}, http.StatusOK, "image/png", "png"},
		{"invalid image type", args{"image/unknown"}, http.StatusUnsupportedMediaType, "", ""},
		{"browser default", args{"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"}, http.StatusOK, "image/png", "png"},
		{"default", args{"image/png"}, http.StatusOK, "image/png", "png"},
		{"default", args{"image/jpg"}, http.StatusOK, "image/jpeg", "jpeg"},
		{"default", args{"image/gif"}, http.StatusOK, "image/gif", "gif"},
		{"default", args{"image/webp"}, http.StatusOK, "image/webp", "webp"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ts := testutils.NewTestServer(ctx, NewAPIv1())

			req := request.Get("%s/api/v1/qrcode", ts.URL).Query("content", "hello world")

			if tt.args.accept != "" {
				req = req.Header(echo.HeaderAccept, tt.args.accept)
			}

			resp, err := req.Do(ctx)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equalf(t, tt.wantStatus, resp.StatusCode, "status not equals: want=%v, got=%v", tt.wantStatus, resp.StatusCode)
			if !resp.Success() {
				return
			}

			require.Equal(t, tt.wantContentType, resp.Header.Get(request.HeaderContentType))

			img, s, err := image.Decode(resp.Body)

			require.NoError(t, err)
			require.Equal(t, tt.wantImage, s)

			got, err := qrcode.Decode(img)
			require.NoError(t, err)
			require.Equal(t, "hello world", got)
		})
	}
}

func TestSize(t *testing.T) {
	type args struct {
		width  int
		height int
	}
	tests := [...]struct {
		name         string
		args         args
		wantW, wantH int
	}{
		{"overflow height", args{0, 2000}, 200, 200},
		{"overflow width", args{2000, 0}, 200, 200},
		{"underflow width", args{-2000, 0}, 200, 200},
		{"default", args{0, 0}, 200, 200},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			ts := testutils.NewTestServer(ctx, NewAPIv1())

			req := request.Get("%s/api/v1/qrcode", ts.URL).Query("content", "hello world")

			if tt.args.width > 0 {
				req = req.Query("w", strconv.FormatInt(int64(tt.args.width), 10))
			}

			if tt.args.height > 0 {
				req = req.Query("h", strconv.FormatInt(int64(tt.args.height), 10))
			}

			resp, err := req.Do(ctx)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			defer resp.Body.Close()
			img, _, err := image.Decode(resp.Body)
			require.NoError(t, err)

			require.NoError(t, err)
			require.Equal(t, image.Point{tt.wantW, tt.wantH}, img.Bounds().Size(), "size not equals")
		})
	}
}

func TestURL(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := testutils.NewTestServer(ctx, NewAPIv1())

	resp, err := request.Get("%s/api/v1/qrcode", ts.URL).
		Query("url", "google.com").Do(ctx)
	require.NoError(t, err)
	require.True(t, resp.Success())

	require.Equal(t, "image/png", resp.Header.Get(request.HeaderContentType))

	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	require.NoError(t, err)
	got, err := qrcode.Decode(img)
	require.NoError(t, err)
	require.Equal(t, "URLTO:google.com", got)
}

func TestWifi(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := testutils.NewTestServer(ctx, NewAPIv1())

	type args struct {
		query    map[string]string
		wantCode string
	}
	tests := [...]struct {
		name       string
		arg        args
		wantErr    bool
		wantStatus int
	}{
		{"empty auth", args{map[string]string{
			"ssid": "myssid",
		}, ""}, false, http.StatusBadRequest},
		{"valid", args{map[string]string{
			"ssid":   "myssid",
			"auth":   "WPA",
			"pass":   "mypassword",
			"hidden": "true",
			"eap":    "TTLS",
			"anon":   "anon_id",
			"ident":  "my_ident",
			"ph2":    "MSCHAPV2",
		}, "WIFI:S:myssid;T:WPA;P:mypassword;H:true;E:TTLS;A:anon_id;I:my_ident;PH2:MSCHAPV2;;"}, false, http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := request.Get("%s/api/v1/qrcode", ts.URL).Queries(tt.arg.query).Do(ctx)
			if (err != nil) != tt.wantErr {
				require.Failf(t, `wifi request failed`, `error = %v, wantErr = %v`, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			require.Equalf(t, tt.wantStatus, resp.StatusCode, "status=%d, wantCode=%d", resp.StatusCode, tt.wantStatus)
			if !resp.Success() {
				return
			}

			require.Truef(t, resp.Success(), "failed with status %s", resp.Status)
			require.Equal(t, "image/png", resp.Header.Get(request.HeaderContentType))

			defer resp.Body.Close()
			img, _, err := image.Decode(resp.Body)
			require.NoError(t, err)

			decoded, err := qrcode.Decode(img)
			require.NoError(t, err)
			require.Equal(t, tt.arg.wantCode, decoded)
		})
	}
}

func TestContact(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := testutils.NewTestServer(ctx, NewAPIv1())

	resp, err := request.Get("%s/api/v1/contact", ts.URL).
		Query("name[first]", "firstname").
		Query("name[last]", "lastname").
		Do(ctx)
	require.NoError(t, err)
	require.True(t, resp.Success())

	require.Equal(t, "image/png", resp.Header.Get(request.HeaderContentType))
}

func TestContactVCF(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := testutils.NewTestServer(ctx, NewAPIv1())

	content := `BEGIN:VCARD
VERSION:4.0
N:lastname;firstname;;;
END:VCARD`

	resp, err := request.Post("%s/api/v1/vcard", ts.URL).
		ContentType(mimeVCard).
		Body(strings.NewReader(content)).
		Do(ctx)
	require.NoError(t, err)
	require.True(t, resp.Success(), "failed with status %d: %s", resp.StatusCode, resp.Status)

	require.Equal(t, "image/png", resp.Header.Get(request.HeaderContentType))
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	require.NoError(t, err)
	got, err := qrcode.Decode(img)
	require.NoError(t, err)
	require.Equal(t, strings.ReplaceAll(content, "\n", "\r\n"), got)
}

// VEvent는 QR 스캐너에서 안되네
func TestVEvent(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ts := testutils.NewTestServer(ctx, NewAPIv1())

	content := `BEGIN:VEVENT
SUMMARY:Summer+Vacation!
DTSTART:20180601T070000Z
DTEND:20180831T070000Z
END:VEVENT`

	resp, err := request.Post("%s/api/v1/vevent", ts.URL).
		ContentType(mimeVEvent).
		Body(strings.NewReader(content)).
		Do(ctx)
	require.NoError(t, err)
	require.True(t, resp.Success(), "failed with status %d: %s", resp.StatusCode, resp.Status)
	require.Equal(t, "image/png", resp.Header.Get(request.HeaderContentType))

	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	require.NoError(t, err)
	got, err := qrcode.Decode(img)
	require.NoError(t, err)

	evt := new(ical.VEvent)
	err = ical.NewEventDecoder(strings.NewReader(got)).Decode(evt)
	require.NoError(t, err)

	require.Equal(t, "Summer+Vacation!", evt.Summary)
	require.Equal(t, ical.DateTime{
		Time: time.Date(2018, 6, 1, 7, 0, 0, 0, time.UTC),
	}, evt.DtStart)
	require.Equal(t, ical.DateTime{
		Time: time.Date(2018, 8, 31, 7, 0, 0, 0, time.UTC),
	}, evt.DtEnd)
}
