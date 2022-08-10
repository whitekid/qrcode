package qrcode

import (
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/emersion/go-vcard"
	"github.com/stretchr/testify/require"
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
				HomeEmail:  "email@home",
				WorkEmail:  "email@work",
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
				"EMAIL;type=INTERNET;type=HOME;type=pref": []*vcard.Field{{Value: "email@home"}},
				"EMAIL;type=INTERNET;type=WORK":           []*vcard.Field{{Value: "email@work"}},
				"ADDR;type=HOME;type=pref":                []*vcard.Field{{Value: ";;;;;;korea"}},
			})
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr, err := tt.args.gen()
			require.NoError(t, err)

			img, err := qr.Render(200, 200)
			require.NoError(t, err)

			// create the output file
			file, _ := os.Create(strings.ReplaceAll(t.Name(), "/", "_") + ".png")
			defer file.Close()

			// encode the qrcode as png
			err = png.Encode(file, img)
			require.NoError(t, err)
		})
	}
}

func TestWifiAuth(t *testing.T) {
	require.Equal(t, AuthNone, StrToWifiAuth("xx"))
}
