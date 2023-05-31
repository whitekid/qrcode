package qrcode

import (
	"fmt"
	"image"
	"strings" 

	"github.com/emersion/go-vcard"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/types"
	"github.com/whitekid/goxp/validate"

	"qrcodeapi/pkg/ical"
)

type QR struct {
	Content string
}

func (q *QR) Render(width, height int) (image.Image, error) {
	return qrcode.NewQRCodeWriter().
		Encode(q.Content, gozxing.BarcodeFormat_QR_CODE,
			width, height, nil)
}

func Text(text string) (*QR, error) {
	if err := validate.Struct(&struct {
		Text string `validate:"required,max=1024"`
	}{
		Text: text,
	}); err != nil {
		return nil, err
	}

	log.Debugf("Text: %s", text)
	return &QR{Content: text}, nil
}

type WiFiAuth int

const (
	AuthNone WiFiAuth = iota
	AuthWEP
	AuthWPA
	AuthWPA2
)

var (
	authStrMap = map[WiFiAuth]string{
		AuthNone: "",
		AuthWEP:  "WEP",
		AuthWPA:  "WPA",
		AuthWPA2: "WPA2",
	}
	strToAuthMap = fx.MapItems(authStrMap, func(k WiFiAuth, v string) (string, WiFiAuth) { return v, k })
)

func (e WiFiAuth) String() string     { return authStrMap[e] }
func StrToWifiAuth(s string) WiFiAuth { return strToAuthMap[s] }

type WPA2Options struct {
	EAPMethod         string `validate:"max=20"`
	AnonymousIdentity string `validate:"max=20"`
	Identity          string `validate:"max=20"`
	Phase2Method      string `validate:"max=20"`
}

// WIFI generate QRCode for joining wifi network
// enc: WEP|WPA|blank
func WIFI(ssid string, auth WiFiAuth, password string, hidden *bool, wpa2 WPA2Options) (*QR, error) {
	if err := validate.Struct(&struct {
		SSID        *string      `validate:"max=20"`
		Password    *string      `validate:"max=20"`
		Hidden      *bool        `validate:"omitempty"`
		WPA2Options *WPA2Options `validate:"omitempty,dive"`
	}{
		SSID:        &ssid,
		Password:    &password,
		Hidden:      hidden,
		WPA2Options: &wpa2,
	}); err != nil {
		return nil, err
	}

	var hiddenStr string

	if hidden != nil {
		hiddenStr = goxp.Ternary(*hidden, "true", "false")
	}

	values := types.NewOrderedMap[string, string]()
	values.Set("S", ssid)
	values.Set("T", auth.String())
	values.Set("P", password)
	values.Set("H", hiddenStr)
	values.Set("E", wpa2.EAPMethod)
	values.Set("A", wpa2.AnonymousIdentity)
	values.Set("I", wpa2.Identity)
	values.Set("PH2", wpa2.Phase2Method)

	values2 := []string{}
	values.ForEach(func(_ int, k string, v string) bool {
		if v == "" {
			return true
		}

		spec := `\;,":`
		for _, x := range spec {
			v = strings.ReplaceAll(v, string(x), `\`+string(x))
		}

		values2 = append(values2, k+":"+v)
		return true
	})

	return Text("WIFI:" + strings.Join(values2, ";") + ";;")
}

type Card struct {
	LastName      string `validate:"max=100"`
	FirstName     string `validate:"max=100"`
	MiddleName    string `validate:"max=100"`
	PrefixName    string `validate:"max=100"` // 경칭
	SuffixName    string `validate:"max=100"` // 호칭
	FormattedName string `validate:"max=100"` // vcf
	NickName      string `validate:"max=100"`

	Company    string `validate:"max=100"` // 회사
	Department string `validate:"max=100"` // 부서
	JobTitle   string `validate:"max=100"` // vcf: 직책

	Email     string `validate:"omitempty,email,max=100"`
	HomeEmail string `validate:"omitempty,email,max=100"`
	WorkEmail string `validate:"omitempty,email,max=100"`

	Tel     string `validate:"max=100"` // MAIN
	Mobile  string `validate:"max=100"`
	HomeTel string `validate:"max=100"`
	WorkTel string `validate:"max=100"`

	HomeFax string `validate:"max=100"`
	WorkFax string `validate:"max=100"`

	Pager string `validate:"max=100"`

	HomeAddr Address `validate:"dive"`
	WorkAddr Address `validate:"dive"`

	Homepage     string `validate:"omitempty,url,max=100"`
	WorkHomepage string `validate:"omitempty,url,max=100"`
	HomeHomepage string `validate:"omitempty,url,max=100"`

	Note string `validate:"max=1024"`

	SocialProfiles []SocialProfile `validate:"dive"`
}

type Address struct {
	PostCode        string `validate:"max=100"`
	CountryOrRegion string `validate:"max=100"`
	Province        string `validate:"max=100"`
	City            string `validate:"max=100"`
	Street          string `validate:"max=100"`
	Street2         string `validate:"max=100"`
}

func (addr *Address) String() string {
	if addr.PostCode == "" && addr.CountryOrRegion == "" && addr.Province == "" && addr.City == "" && addr.Street == "" && addr.Street2 == "" {
		return ""
	}

	street := goxp.Ternary(addr.Street2 == "", addr.Street, addr.Street+"\n"+addr.Street2)
	return strings.Join([]string{"", "", street, addr.City, addr.Province, addr.PostCode, addr.CountryOrRegion}, ";")
}

type SocialProfile struct {
	Type string `validate:"max=100"`
	ID   string `validate:"max=100"`
}

func setField(vc vcard.Card, key, value string) {
	_, exists := vc[key]

	if (exists && value == "") || value != "" {
		vc[key] = []*vcard.Field{{Value: value}}
	}
}

func Contact(card *Card) (*QR, error) {
	if err := validate.Struct(card); err != nil {
		return nil, err
	}

	vc := vcard.Card{
		"VERSION": []*vcard.Field{{Value: "4.0"}},
		"N":       []*vcard.Field{{Value: fmt.Sprintf("%s;%s;%s;%s;%s", card.LastName, card.FirstName, card.MiddleName, card.PrefixName, card.SuffixName)}}, // Name
	}

	setField(vc, "EMAIL;type=INTERNET;type=HOME;type=pref", card.HomeEmail)
	setField(vc, "EMAIL;type=INTERNET;type=WORK", card.WorkEmail)
	setField(vc, "TEL;type=CELL;type=VOICE;type=pref", card.Mobile)
	setField(vc, "TEL;type=HOME;type=VOICE", card.HomeTel)
	setField(vc, "TEL;type=WORK;type=VOICE", card.WorkTel)
	setField(vc, "TEL;type=MAIN", card.Tel)

	setField(vc, "TEL;type=HOME;type=FAX", card.HomeFax)
	setField(vc, "TEL;type=WORK;type=FAX", card.WorkFax)

	setField(vc, "TEL;type=PAGER", card.Pager)
	setField(vc, "ADR;type=HOME;type=pref", card.HomeAddr.String())
	setField(vc, "ADR;type=WORK", card.WorkAddr.String())

	for _, social := range card.SocialProfiles {
		typ := social.Type
		ID := social.ID

		switch social.Type {
		case "twitter":
			ID = "http://twitter.com/" + ID
		case "facebook":
			ID = "http://facebook.com/" + ID
		case "flickr":
			ID = "http://www.fickr.com/photos/" + ID
		case "linkedin":
			ID = "http://www.fickr.com/in/" + ID
		case "myspace":
			ID = "http://www.myspace.com/" + ID
		case "sinaweibo":
			ID = "http://weibo.com/n/" + ID
		case "JabberInstant":
			typ = "X-SOCIALPROFILE;type=JabberInstant;x-user=" + ID
			ID = "http://t.qq.com/" + ID
		case "yelp":
			typ = "X-SOCIALPROFILE;type=Yelp:x-apple"
		default:
			return nil, fmt.Errorf("unsupported social type: %s", social.Type)
		}
		setField(vc, "X-SOCIALPROFILE;type="+typ, ID)
	}

	setField(vc, "NOTE", card.Note)

	s := new(strings.Builder)
	if err := vcard.NewEncoder(s).Encode(vc); err != nil {
		return nil, err
	}

	return Text(s.String())
}

func VCard(card vcard.Card) (*QR, error) {
	s := new(strings.Builder)

	if err := vcard.NewEncoder(s).Encode(card); err != nil {
		return nil, err
	}

	return Text(strings.TrimSpace(s.String()))
}

type ICSVEvent struct {
}

func VEvent(evt *ical.VEvent) (*QR, error) {
	buf := new(strings.Builder)
	if err := ical.NewEventEncoder(buf).Encode(evt); err != nil {
		return nil, err
	}

	return Text(buf.String())
}
