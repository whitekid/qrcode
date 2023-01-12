package echox

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"qrcodeapi/config"
)

type Echo struct {
	*echo.Echo
}

func New() *Echo {
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
	e.Use(middleware.CORS())

	return &Echo{e}
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

type Router interface {
	Route(e *Echo, path string)
}
