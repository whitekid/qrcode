package qrcode

import (
	"bytes"
	"fmt"
	"image"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/types"
)

type QR struct {
	Content string
}

func (q *QR) Render(width, height int) (image.Image, error) {
	return qrcode.NewQRCodeWriter().
		Encode(q.Content, gozxing.BarcodeFormat_QR_CODE,
			width, height, nil)
}

func Text(content string) (*QR, error) { return &QR{Content: content}, nil }

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
	strToAuthMap = map[string]WiFiAuth{
		"WEP":  AuthWEP,
		"WPA":  AuthWPA,
		"WPA2": AuthWPA2,
	}
)

func (e WiFiAuth) String() string     { return authStrMap[e] }
func StrToWifiAuth(s string) WiFiAuth { return strToAuthMap[s] }

type WPA2Options struct {
	EAPMethod         string
	AnonymousIdentity string
	Identity          string
	Phase2Method      string
}

// WIFI generate QRCode for joining wifi network
// enc: WEP|WPA|blank
func WIFI(ssid string, auth WiFiAuth, password string, hidden *bool,
	wpa2 WPA2Options) (*QR, error) {
	var hiddenStr string

	if hidden != nil {
		hiddenStr = fx.Ternary(*hidden, "true", "false")
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
	LastName      string
	FirstName     string
	MiddleName    string
	PrefixName    string // 경칭
	SuffixName    string // 호칭
	FormattedName string // vcf
	NickName      string

	Company    string // 회사
	Department string // 부서
	JobTitle   string // vcf: 직책

	Email     string
	HomeEmail string
	WorkEmail string

	Tel     string // MAIN
	Mobile  string
	HomeTel string
	WorkTel string

	HomeFax string
	WorkFax string

	Pager string

	HomeAddr Address
	WorkAddr Address

	Homepage     string
	WorkHomepage string
	HomeHomepage string

	Note string

	SocialProfiles []SocialProfile
}

type Address struct {
	PostCode        string
	CountryOrRegion string
	Province        string
	City            string
	Street          string
	Street2         string
}

func (addr *Address) String() string {
	if addr.PostCode == "" && addr.CountryOrRegion == "" && addr.Province == "" && addr.City == "" && addr.Street == "" && addr.Street2 == "" {
		return ""
	}

	street := fx.Ternary(addr.Street2 == "", addr.Street, addr.Street+"\n"+addr.Street2)
	return strings.Join([]string{"", "", street, addr.City, addr.Province, addr.PostCode, addr.CountryOrRegion}, ";")
}

type SocialProfile struct {
	Type string
	ID   string
}

func setField(vc vcard.Card, key, value string) {
	_, exists := vc[key]

	if (exists && value == "") || value != "" {
		vc[key] = []*vcard.Field{{Value: value}}
	}
}

func Contact(card *Card) (*QR, error) {
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

	var buf bytes.Buffer
	if err := vcard.NewEncoder(&buf).Encode(vc); err != nil {
		return nil, err
	}

	return Text(buf.String())
}

func VCard(card vcard.Card) (*QR, error) {
	var buf bytes.Buffer

	if err := vcard.NewEncoder(&buf).Encode(card); err != nil {
		return nil, err
	}

	return Text(strings.TrimSpace(buf.String()))
}
