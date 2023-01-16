package qrcode

import (
	"bytes"
	"image/png"
	"strings"
	"testing"
	"time"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/require"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/fx"

	"qrcodeapi/pkg/ical"
)

func TestQR(t *testing.T) {
	True := true
	False := false

	type args struct {
		gen func() (*QR, error)
	}
	tests := [...]struct {
		name string
		args args
	}{
		{"text", args{func() (*QR, error) { return Text("hello") }}},
		{"url", args{func() (*QR, error) { return Text("URLTO:http://google.com") }}},
		{"email", args{func() (*QR, error) {
			return Text("mailto:email@host.com?subject=subject&cc=cc1@host.com,cc2@host.com&bcc=bcc@host.com@body=hello")
		}}},
		{"tel", args{func() (*QR, error) { return Text("tel:+123456789") }}},
		{"sms", args{func() (*QR, error) { return Text("sms:+123456789:message%20here.") }}},
		{"facetime", args{func() (*QR, error) { return Text("facetime:+123456789") }}},
		{"facetime", args{func() (*QR, error) { return Text("facetime:me@icloud.com") }}},
		{"facetime-audio", args{func() (*QR, error) { return Text("facetime:me@icloud.com") }}},
		{"playstore", args{func() (*QR, error) { return Text("market://details?id=org.example.foo") }}},
		{"wifi", args{func() (*QR, error) { return WIFI("SSID", AuthWPA, "", nil, WPA2Options{}) }}},
		{"wifi-hidden", args{func() (*QR, error) {
			return WIFI("SSID", AuthWPA, "", &True, WPA2Options{})
		}}},
		{"wifi-not-hidden", args{func() (*QR, error) {
			return WIFI("SSID", AuthWPA, "", &False, WPA2Options{})
		}}},
		{"contact", args{func() (*QR, error) {
			return Contact(&Card{
				FirstName:  "firstName",
				LastName:   "lastName",
				MiddleName: "middleName",
				PrefixName: "prefix",
				SuffixName: "suffix",
				HomeEmail:  "email@home.com",
				WorkEmail:  "email@work.com",
				HomeAddr:   Address{CountryOrRegion: "korea"},
				WorkAddr:   Address{Province: "at work"},
				SocialProfiles: []SocialProfile{
					{Type: "twitter", ID: "twitter_id"},
				},
			})
		}}},
		{"vcard", args{func() (*QR, error) {
			return VCard(vcard.Card{
				"VERSION": []*vcard.Field{{Value: "4.0"}},
				"N":       []*vcard.Field{{Value: "lastName;firstName;middleName;prefix;suffix"}},
				"EMAIL;type=INTERNET;type=HOME;type=pref": []*vcard.Field{{Value: "email@home.com"}},
				"EMAIL;type=INTERNET;type=WORK":           []*vcard.Field{{Value: "email@work.com"}},
				"ADDR;type=HOME;type=pref":                []*vcard.Field{{Value: ";;;;;;korea"}},
			})
		}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qr, err := tt.args.gen()
			require.NoError(t, err)

			img, err := qr.Render(200, 200)
			require.NoError(t, err)

			// encode the qrcode as png
			err = png.Encode(new(bytes.Buffer), img)
			require.NoError(t, err)
		})
	}
}

func TestText(t *testing.T) {
	type args struct {
		text string
	}
	tests := [...]struct {
		name    string
		args    args
		textErr bool
		wantErr bool
	}{
		{`empty`, args{""}, true, false},
		{`valid`, args{"동해물과"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := Text(tt.args.text)
			require.Truef(t, (err != nil) == tt.textErr, `Text() failed: error = %+v, textErr = %v`, err, tt.textErr)
			if tt.textErr {
				return
			}

			got, err := q.Render(200, 200)
			require.Truef(t, (err != nil) == tt.wantErr, `Render() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}

			dec, err := Decode(got)
			require.NoError(t, err)
			require.Equal(t, tt.args.text, dec)
		})
	}
}

func FuzzText(f *testing.F) {
	f.Add("동해물과")
	f.Fuzz(func(t *testing.T, text string) {
		if len(text) == 0 || len(text) > 1024 {
			t.Skip()
		}

		q, err := Text(text)
		require.NoError(t, err)

		_, err = q.Render(200, 200)
		require.NoErrorf(t, err, "text: %d", len(text))
	})
}

func BenchmarkText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		text := goxp.RandomString(30)
		q, err := Text(text)
		require.NoError(b, err)

		img, err := q.Render(200, 200)
		require.NoError(b, err)

		_ = img
	}
}

func TestWifi(t *testing.T) {
	type args struct {
		ssid     string
		auth     WiFiAuth
		password string
	}
	tests := [...]struct {
		name      string
		args      args
		wifiErr   bool
		renderErr bool
	}{
		{`empty`, args{}, false, false},
		{`long ssid`, args{ssid: "012345678990123456789011"}, true, false},
		{`long password`, args{password: "012345678990123456789011"}, true, false},
		{`valid`, args{"ssid", AuthNone, "pass"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := WIFI(tt.args.ssid, tt.args.auth, tt.args.password, nil, WPA2Options{})
			require.Truef(t, (err != nil) == tt.wifiErr, `Wifi() failed: error = %+v, wifiErr = %v`, err, tt.wifiErr)
			if tt.wifiErr {
				return
			}

			got, err := q.Render(200, 200)
			require.Truef(t, (err != nil) == tt.renderErr, `Render() failed: error = %+v, renderErr = %v`, err, tt.renderErr)
			if tt.renderErr {
				return
			}
			_ = got
		})
	}
}

// TODO move to goxp
func NoError1[T1 any](t *testing.T, v fx.Tuple2[T1, error]) T1 {
	require.NoError(t, v.V2)
	return v.V1
}

func FuzzWifi(f *testing.F) {
	f.Add("ssid", "password")
	f.Fuzz(func(t *testing.T, ssid string, password string) {
		if len(ssid) > 20 || len(password) > 20 {
			t.Skip()
		}

		q := NoError1(t, fx.T2(WIFI(ssid, AuthNone, password, nil, WPA2Options{})))
		NoError1(t, fx.T2((q.Render(200, 200))))
	})
}

func TestContact(t *testing.T) {
	card := &Card{}
	qr, err := Contact(card)
	require.NoError(t, err)

	img, err := qr.Render(200, 200)
	require.NoError(t, err)

	s, err := Decode(img)
	require.NoError(t, err)
	require.Regexp(t, `^BEGIN:VCARD`, s)
}

func FuzzContact(f *testing.F) {
	f.Add("Last Name", "middle name", "first name")
	f.Fuzz(func(t *testing.T, lastName, middleName, firstName string) {
		if len(lastName) > 100 || len(middleName) > 100 || len(firstName) > 100 {
			t.Skip()
		}

		card := &Card{
			LastName:   lastName,
			MiddleName: middleName,
			FirstName:  firstName,
		}
		qr, err := Contact(card)
		require.NoError(t, err)

		img, err := qr.Render(200, 200)
		require.NoError(t, err)

		_ = img
	})
}

func TestWifiAuth(t *testing.T) {
	require.Equal(t, AuthNone, StrToWifiAuth("xx"))
}

func TestVEvent(t *testing.T) {
	type args struct {
		event *ical.VEvent
	}
	tests := [...]struct {
		name    string
		args    args
		wantErr bool
	}{
		{`valid`, args{&ical.VEvent{
			Summary:     "summary",
			Description: "동해물과 백두산이 마르고 닳도록 동해물과 백두산이 마르고 닳도록 동해물과 백두산이 마르고 닳도록 동해물과 백두산이 마르고 닳도록",
			DtStamp:     ical.DateTime{Time: time.Date(2023, 1, 16, 3, 4, 5, 0, time.UTC)},
		}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr, err := VEvent(tt.args.event)
			require.Truef(t, (err != nil) == tt.wantErr, `VEvent() failed: error = %+v, wantErr = %v`, err, tt.wantErr)

			img, err := qr.Render(200, 200)
			require.NoError(t, err)

			s, err := Decode(img)
			require.NoError(t, err)
			require.Regexp(t, `^BEGIN:VEVENT`, s)

			evt := new(ical.VEvent)
			require.NoError(t, ical.NewEventDecoder(strings.NewReader(s)).Decode(evt))
			require.Equal(t, tt.args.event, evt)
		})
	}
}

func FuzzVEvent(f *testing.F) {
	f.Add("summary", "description", time.Now().Unix())
	f.Fuzz(func(t *testing.T, summary, description string, dtStamp int64) {
		evt := &ical.VEvent{
			Summary:     summary,
			Description: description,
			DtStamp:     ical.DateTime{Time: time.Unix(dtStamp, 0)},
		}
		qr, err := VEvent(evt)
		require.NoError(t, err)

		img, err := qr.Render(200, 200)
		require.NoError(t, err)

		s, err := Decode(img)
		require.NoError(t, err)
		require.Regexp(t, `^BEGIN:VEVENT`, s)
	})
}
