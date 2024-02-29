package main

import (
	"expvar"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) routes() http.Handler {

	e := echo.New()

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(app.metrics)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiterWithConfig(config))
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.Use(middleware.BodyLimit("2K"))

	e.GET("/v1/ping", app.ping)

	// User Routes

	e.POST("/v1/admin", app.registerAdminHandler)
	e.PUT("/v1/admin/activate", app.activateAdminHandler)
	e.POST("/v1/login", app.createAuthenticationTokenHandler)
	e.POST("/v1/token/resend/activation", app.createActivationTokenHandler)

	// password management
	e.POST("/v1/token/password/reset", app.createPasswordResetTokenHandler)
	e.PUT("/v1/admin/password", app.updateAdminPasswordHandler)

	// metrics
	e.GET("/v1/metrics", echo.WrapHandler(expvar.Handler()))

	g := e.Group("/v1/auth")

	g.Use(app.authenticate)

	// houses
	g.GET("/house", app.listHousesHandler, app.requireAuthenticatedAdmin)
	g.POST("/house", app.createHouseHandler, app.requireAuthenticatedAdmin)
	g.POST("/bulk/house", app.bulkHousesHandler, app.requireAuthenticatedAdmin)
	g.GET("/house/:uuid", app.showHousesHandler, app.requireAuthenticatedAdmin)
	g.PUT("/house/:uuid", app.updateHouseHandler, app.requireAuthenticatedAdmin)

	// tenants
	g.GET("/tenant", app.listTenantsHandler, app.requireAuthenticatedAdmin)
	g.POST("/tenant", app.createTenantHandler, app.requireAuthenticatedAdmin)
	g.GET("/tenant/:uuid", app.showTenantsHandler, app.requireAuthenticatedAdmin)
	g.PUT("/tenant/:uuid", app.updateTenantsHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/tenant/:uuid", app.removeTenant, app.requireAuthenticatedAdmin)

	// Payments
	g.GET("/payment", app.listPaymentsHandler, app.requireAuthenticatedAdmin)
	g.POST("/payment", app.createPaymentHandler, app.requireAuthenticatedAdmin)
	g.GET("/payment/:uuid", app.showPaymentHandler, app.requireAuthenticatedAdmin)
	g.PUT("/payment/:uuid", app.updatePaymentHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/payment/:uuid", app.deletePaymentHandler, app.requireAuthenticatedAdmin)

	return e

}
