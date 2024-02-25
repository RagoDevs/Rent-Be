package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) routes() http.Handler {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))

	e.GET("/v1/ping", app.ping)

	// User Routes

	e.POST("/v1/admin", app.registerAdminHandler)
	e.PUT("/v1/admin/activate", app.activateAdminHandler)
	e.POST("/v1/login", app.createAuthenticationTokenHandler)
	e.POST("/v1/token/resend/activation", app.createActivationTokenHandler)

	// password management
	e.POST("/v1/token/password/reset", app.createPasswordResetTokenHandler)
	e.PUT("/v1/admin/password", app.updateAdminPasswordHandler)

	g := e.Group("/v1/auth")

	g.Use(app.authenticate)

	// houses
	g.GET("/house", app.listHousesHandler, app.requireAuthenticatedAdmin)
	g.POST("/house", app.createHouseHandler, app.requireAuthenticatedAdmin)
	g.POST("/bulk/houses", app.bulkHousesHandler)
	g.GET("/house/:uuid", app.showHousesHandler, app.requireAuthenticatedAdmin)
	g.PUT("/house/:uuid", app.updateHouseHandler, app.requireAuthenticatedAdmin)

	// tenants
	g.GET("/v1/tenant", app.listTenantsHandler, app.requireAuthenticatedAdmin)
	g.POST("/v1/tenant", app.createTenantHandler, app.requireAuthenticatedAdmin)
	g.GET("/v1/tenant/:uuid", app.showTenantsHandler, app.requireAuthenticatedAdmin)
	g.PUT("/v1/tenant/:uuid", app.updateTenantsHandler, app.requireAuthenticatedAdmin)
	g.DELETE("/v1/tenant/:uuid", app.removeTenant, app.requireAuthenticatedAdmin)

	// Payments
	g.GET("/v1/payment", app.listPaymentsHandler, app.requireAuthenticatedAdmin)
	g.POST("/v1/payment", app.createPaymentHandler, app.requireAuthenticatedAdmin)
	g.GET("/v1/payment/:uuid", app.showPaymentHandler, app.requireAuthenticatedAdmin)
	g.PUT("/v1/payment/:uuid", app.updatPaymentHandler, app.requireAuthenticatedAdmin)

	return e

}
