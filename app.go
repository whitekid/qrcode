package qrcodeapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/whitekid/goxp/fx"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/service"
	"golang.org/x/time/rate"

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

	if err := e.Start(config.BindAddr()); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
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

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Validator = &Validator{validator: validator.New()}
	e.Use(func(logCode int) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			// log http errors
			return func(c echo.Context) error {
				err := next(c)
				if err != nil {
					code := http.StatusInternalServerError

					if ee, ok := err.(validator.ValidationErrors); ok {
						err = echo.NewHTTPError(http.StatusBadRequest, ee.Error())
					}

					if he, ok := err.(*echo.HTTPError); ok {
						code = he.Code
					}

					if code >= logCode {
						c.Logger().Errorf("%+v", err)
					}
				}

				return err
			}
		}
	}(http.StatusBadRequest))

	e.Use(middleware.Logger())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(config.RateLimit()))))

	return e
}

func (s *qrcodeService) setup() *echo.Echo {
	e := newEcho()
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "https://github.com/whitekid/qrcodeapi")
	})
	newAPIv1().Route(e, "/v1")

	return e
}

func parseIntDef(s string, defaultValue, minValue, maxValue int) int {
	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return fx.Min([]int{fx.Max([]int{value, minValue}), maxValue})
}
