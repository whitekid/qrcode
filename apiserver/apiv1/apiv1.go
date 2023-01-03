package apiv1

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/labstack/echo/v4"
	"github.com/whitekid/goxp/request"

	"qrcodeapi/pkg/helper"
	"qrcodeapi/pkg/helper/echox"
	"qrcodeapi/pkg/qrcode"
)

type APIv1 struct{}

var _ echox.Router = (*APIv1)(nil)

func NewAPIv1() echox.Router { return &APIv1{} }

func (api *APIv1) Route(e *echox.Echo, path string) {
	v1 := e.Group(path)

	v1.GET("/qrcode", api.handleGenerate)
	v1.GET("/contact", api.handleContact)
	v1.POST("/vcard", api.handleContactVCard)
	v1.POST("/vevent", api.handleVEvent)
}

type RenderRequest struct {
	W         int    `query:"w"`
	H         int    `query:"h"`
	ImageType string `header:"accept"`
}

func (api *APIv1) renderQRCode(c echo.Context, in *qrcode.QR) error {
	// NOTE c.Bind()는 Post에서 동작하지 않음
	req := &RenderRequest{
		W:         helper.ParseIntDef(c.QueryParam("w"), 200, 21, 200),
		H:         helper.ParseIntDef(c.QueryParam("h"), 200, 21, 200),
		ImageType: c.Request().Header.Get(echo.HeaderAccept),
	}

	img, err := in.Render(req.W, req.H)
	if err != nil {
		return err
	}

	switch strings.TrimPrefix(strings.ToLower(req.ImageType), "image/") {
	case "jpeg", "jpg":
		return jpeg.Encode(c.Response().Writer, img, nil)
	case "gif":
		return gif.Encode(c.Response().Writer, img, nil)
	case "", "png":
		return png.Encode(c.Response().Writer, img)
	default:
		return echo.ErrUnsupportedMediaType
	}
}

type GenerateRequest struct {
	Content string `query:"content"`
	URL     string `query:"url"`
	SSID    string `query:"ssid"`
}

func (api *APIv1) handleGenerate(c echo.Context) error {
	req := &GenerateRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	switch {
	case req.Content != "":
		qr, err := qrcode.Text(req.Content)
		if err != nil {
			return err
		}
		return api.renderQRCode(c, qr)

	case req.URL != "":
		qr, err := qrcode.Text("URLTO:" + req.URL)
		if err != nil {
			return err
		}
		return api.renderQRCode(c, qr)

	case req.SSID != "":
		return api.handleWifi(c)
	}

	return echo.NewHTTPError(http.StatusBadRequest)
}

type WIFIRequest struct {
	SSID   string `query:"ssid" validate:"required"`
	Auth   string `query:"auth" validate:"required"`
	Pass   string `query:"pass"`
	Hidden string `query:"hidden"`
	EAP    string `query:"eap"`
	AnonID string `query:"anon"`
	Ident  string `query:"ident"`
	PH2    string `query:"ph2"`
}

func (api *APIv1) handleWifi(c echo.Context) error {
	req := &WIFIRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	hidden := (*bool)(nil)
	switch req.Hidden {
	case "true":
		v := true
		hidden = &v
	case "false", "":
		v := false
		hidden = &v
	}

	qr, err := qrcode.WIFI(req.SSID,
		qrcode.StrToWifiAuth(req.Auth), req.Pass, hidden,
		qrcode.WPA2Options{
			EAPMethod:         req.EAP,
			AnonymousIdentity: req.AnonID,
			Identity:          req.Ident,
			Phase2Method:      req.PH2})
	if err != nil {
		return err
	}

	return api.renderQRCode(c, qr)
}

type ContactRequest struct {
	FirstName  string `query:"name[first]"`
	LastName   string `query:"name[last]"`
	MiddleName string `query:"name[middle]"`

	Company    string `query:"company"`
	Department string `query:"department"`
	JobTitle   string `query:"title"`

	Email     string `query:"email"`
	EmailHome string `query:"email[home]"`
	EmailWork string `query:"email[work]"`

	Tel     string `query:"tel"`
	TelHome string `query:"tel[home]"`
	TelWork string `query:"tel[work]"`
	Mobile  string `query:"mobile"`
	Pager   string `query:"pager"`

	FaxHome string `query:"fax[home]"`
	FaxWork string `query:"fax[work]"`

	HomeAddr struct {
		PostCode        string `query:"addr[home][postcode]"`
		CountryOrRegion string `query:"addr[home][country]"`
		Province        string `query:"addr[home][province]"`
		City            string `query:"addr[home][city]"`
		Street          string `query:"addr[home][street]"`
		Street2         string `query:"addr[home][street2]"`
	} `validate:"dive"`
	WorkAddr struct {
		PostCode        string `query:"addr[work][postcode]"`
		CountryOrRegion string `query:"addr[work][country]"`
		Province        string `query:"addr[work][province]"`
		City            string `query:"addr[work][city]"`
		Street          string `query:"addr[work][street]"`
		Street2         string `query:"addr[work][street2]"`
	} `validate:"dive"`

	Note string `query:"note"`
}

func (api *APIv1) handleContact(c echo.Context) error {
	req := &ContactRequest{}
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	qr, err := qrcode.Contact(&qrcode.Card{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		MiddleName: req.MiddleName,

		Company:    req.Company,
		Department: req.Department,
		JobTitle:   req.JobTitle,

		Email:     req.Email,
		HomeEmail: req.EmailHome,
		WorkEmail: req.EmailWork,

		Tel:     req.Tel,
		Mobile:  req.Mobile,
		HomeTel: req.TelHome,
		WorkTel: req.TelWork,

		HomeFax: req.FaxHome,
		WorkFax: req.FaxWork,
		Pager:   req.Pager,

		HomeAddr: qrcode.Address{
			PostCode:        req.HomeAddr.PostCode,
			CountryOrRegion: req.HomeAddr.CountryOrRegion,
			Province:        req.HomeAddr.Province,
			City:            req.HomeAddr.City,
			Street:          req.HomeAddr.Street,
			Street2:         req.HomeAddr.Street2,
		},
		WorkAddr: qrcode.Address{
			PostCode:        req.WorkAddr.PostCode,
			CountryOrRegion: req.WorkAddr.CountryOrRegion,
			Province:        req.WorkAddr.Province,
			City:            req.WorkAddr.City,
			Street:          req.WorkAddr.Street,
			Street2:         req.WorkAddr.Street2,
		},

		Note: req.Note,
	})
	if err != nil {
		return err
	}

	return api.renderQRCode(c, qr)
}

const (
	mimeVCard  = "text/vcard"
	mimeVEvent = "text/vevent"
)

func (api *APIv1) handleContactVCard(c echo.Context) error {
	if mediaType, _, _ := mime.ParseMediaType(c.Request().Header.Get(request.HeaderContentType)); mediaType != mimeVCard {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	card, err := vcard.NewDecoder(c.Request().Body).Decode()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	qr, err := qrcode.VCard(card)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return api.renderQRCode(c, qr)
}

func (api *APIv1) handleVEvent(c echo.Context) error {
	if mediaType, _, _ := mime.ParseMediaType(c.Request().Header.Get(request.HeaderContentType)); mediaType != mimeVEvent {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	body, err := io.ReadAll(c.Request().Body)
	defer c.Request().Body.Close()
	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(body), "\n", "\r\n")
	qr, err := qrcode.Text(content)
	if err != nil {
		return err
	}

	return api.renderQRCode(c, qr)
}
