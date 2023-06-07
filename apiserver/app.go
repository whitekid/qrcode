package apiserver

import (
	"context"
	"net/http"

	openapimiddleware "github.com/go-openapi/runtime/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/whitekid/echox"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/service"
	"golang.org/x/time/rate"

	"qrcodeapi/apiserver/apiv1"
	"qrcodeapi/config"
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
	e.File("/swagger.json", "spec/tsp-output/@typespec/openapi3/openapi.json")

	return http.ListenAndServe(config.BindAddr(), h)
}

func (s *qrcodeService) setup() *echox.Echo {
	e := echox.New(
		middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.RateLimit()))),
		middleware.CORS(),
	)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "https://github.com/whitekid/qrcode")
	})
	e.Route(nil, apiv1.NewAPIv1())

	for _, r := range e.Routes() {
		log.Debugf("%s %s => %s", r.Method, r.Path, r.Name)
	}

	return e
}
