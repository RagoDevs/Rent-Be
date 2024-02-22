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

	e.POST("/v1/admin", app.registerUserHandler)
	e.PUT("/v1/admin/activate", app.activateUserHandler)
	e.POST("/v1/login", app.createAuthenticationTokenHandler)
	e.POST("/v1/token/activation", app.createActivationTokenHandler)


	// password management
	e.GET("/v1/token/passwordreset", app.createPasswordResetTokenHandler)
	e.PUT("/v1/admin/password", app.updateUserPasswordHandler)

	g := e.Group("/v1/auth")

	g.Use(app.authenticate)



	// houses
	g.GET("/house", app.listHousesHandler)
	g.POST("/house", app.createHouseHandler)
	// g.POST("/bulk/houses", app.bulkHousesHandler)
	g.GET("/house/:uuid", app.showHousesHandler)
	g.PUT("/house/:uuid", app.updateHouseHandler)

	// tenants
	g.GET("/v1/tenant", app.listTenantsHandler)
	g.POST("/v1/tenant", app.createTenantHandler)
	g.GET("/v1/tenant/:uuid", app.showTenantsHandler)
	g.PUT("/v1/tenant/:uuid", app.updateTenantsHandler)
	g.DELETE("/v1/tenant/:uuid", app.removeTenant)

	// Payments
	g.GET("/v1/payment", app.listPaymentsHandler)
	g.POST("/v1/payment", app.createPaymentHandler)
	g.GET("/v1/payment/:uuid", app.showPaymentHandler)
	g.PUT("/v1/payment/:uuid", app.updatPaymentHandler)


	return e

}
