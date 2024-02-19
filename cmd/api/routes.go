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

	e.POST("/v1/users", app.registerUserHandler)
	e.PUT("/v1/users/activated", app.activateUserHandler)
	e.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	e.POST("/v1/tokens/activation", app.createActivationTokenHandler)


	g := e.Group("/v1/auth")

	g.Use(app.authenticate)

	// User Routes

	g.POST("/v1/users", app.registerUserHandler)
	g.PUT("/v1/users/activated", app.activateUserHandler)
	g.POST("/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	g.POST("/v1/tokens/activation", app.createActivationTokenHandler)

	// houses
	g.GET("/houses", app.listHousesHandler)
	g.POST("/houses", app.createHouseHandler)
	// g.POST("/bulk/houses", app.bulkHousesHandler)
	g.GET("/houses/:uuid", app.showHousesHandler)
	g.PUT("/houses/:uuid", app.updateHouseHandler)

	// tenants
	g.GET("/v1/tenants", app.listTenantsHandler)
	g.POST("/v1/tenants", app.createTenantHandler)
	g.GET("/v1/tenants/:uuid", app.showTenantsHandler)
	g.PUT("/v1/tenants/:uuid", app.updateTenantsHandler)
	g.DELETE("/v1/tenants/:uuid", app.removeTenant)

	// Payments
	g.GET("/v1/payments", app.listPaymentsHandler)
	g.POST("/v1/payments", app.createPaymentHandler)
	g.GET("/v1/payments/:uuid", app.showPaymentHandler)
	g.PUT("/v1/payments/:uuid", app.updatPaymentHandler)

	//password management
	g.GET("/v1/tokens/passwordreset", app.createPasswordResetTokenHandler)
	g.PUT("/v1/users/password", app.updateUserPasswordHandler)

	return e

}
