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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiterWithConfig(config))
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.Use(middleware.BodyLimit("2K"))

	e.GET("/v1/ping", app.ping)

	// User Routes

	e.POST("/v1/admins", app.registerAdminHandler)
	e.PUT("/v1/admins/activate", app.activateAdminHandler)
	e.POST("/v1/login", app.createAuthenticationTokenHandler)
	e.POST("/v1/tokens/resend/activation", app.resendActivationTokenHandler)

	// password management
	e.POST("/v1/tokens/password/reset", app.createPasswordResetTokenHandler)
	e.PUT("/v1/admins/password/reset", app.updateAdminPasswordOnResetHandler)

	// metrics
	e.GET("/v1/metrics", echo.WrapHandler(expvar.Handler()))

	g := e.Group("/v1/auth")

	g.Use(app.authenticate)

	// houses
	g.GET("/houses", app.listHousesHandler, app.requireAuthenticatedAdmin)
	g.POST("/houses", app.createHouseHandler, app.requireAuthenticatedAdmin)
	g.POST("/bulk/houses", app.bulkHousesHandler, app.requireAuthenticatedAdmin)
	g.GET("/houses/:uuid", app.showHousesHandler, app.requireAuthenticatedAdmin)
	g.GET("/houses/tenant/:uuid", app.showHousesTenantHandler, app.requireAuthenticatedAdmin)
	g.PUT("/houses/:uuid", app.updateHouseHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/houses/:uuid", app.deleteHousesHandler, app.requireAuthenticatedAdmin)

	// tenants
	g.GET("/tenants", app.listTenantsHandler, app.requireAuthenticatedAdmin)
	g.POST("/tenants", app.createTenantHandler, app.requireAuthenticatedAdmin)
	g.GET("/tenants/:uuid", app.showTenantsHandler, app.requireAuthenticatedAdmin)
	g.GET("/tenants/house/:uuid", app.showTenantsHouseHandler, app.requireAuthenticatedAdmin)
	g.PUT("/tenants/:uuid", app.updateTenantsHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/tenants/:uuid", app.removeTenant, app.requireAuthenticatedAdmin)

	// Payments
	g.GET("/payments", app.listPaymentsHandler, app.requireAuthenticatedAdmin)
	g.POST("/payments", app.createPaymentHandler, app.requireAuthenticatedAdmin)
	g.GET("/payments/:uuid", app.showPaymentHandler, app.requireAuthenticatedAdmin)
	g.PUT("/payments/:uuid", app.updatePaymentHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/payments/:uuid", app.deletePaymentHandler, app.requireAuthenticatedAdmin)

	return e

}
