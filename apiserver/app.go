package apiserver

import (
	"context"
	"net/http"

	openapimiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/service"

	"qrcodeapi/apiserver/apiv1"
	"qrcodeapi/config"
	"qrcodeapi/pkg/helper/echox"
)

func Run(ctx context.Context) error { return New().Serve(ctx) }

type qrcodeService struct{}

var _ service.Interface = (*qrcodeService)(nil)

func New() service.Interface { return &qrcodeService{} }

func (s *qrcodeService) Serve(ctx context.Context) error {
	e := s.setup()

	go func() {
		<-ctx.Done()
		if err := e.Shutdown(context.Background()); err != nil {
			log.Fatalf("%s", err)
		}
	}()

	// TODO create redoc middleware
	h := openapimiddleware.Redoc(openapimiddleware.RedocOpts{
		Path:    "swagger-ui",
		SpecURL: "/swagger.json",
	}, e)
	e.File("/swagger.json", "cadl/cadl-output/@cadl-lang/openapi3/openapi.json")

	return http.ListenAndServe(config.BindAddr(), h)
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *qrcodeService) setup() *echox.Echo {
	e := echox.New()
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "https://github.com/whitekid/qrcodeapi")
	})
	apiv1.NewAPIv1().Route(e, "/api/v1")

	return e
}
