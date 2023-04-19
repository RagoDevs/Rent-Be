package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/v1/ping", app.ping)

	// test

	router.HandlerFunc(http.MethodGet, "/v1/test", app.requireActivatedUser(app.test))

	// User Routes

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Resend user token route
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)

	// Metrics Routes
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	//password management

	router.HandlerFunc(http.MethodPost, "/v1/tokens/passwordreset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)

	// houses

	router.HandlerFunc(http.MethodGet, "/v1/houses", app.requireActivatedUser(app.listHousesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/houses", app.requireActivatedUser(app.createHouseHandler))
	router.HandlerFunc(http.MethodPost, "/v1/bulk/houses", app.requireActivatedUser(app.bulkHousesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/houses/:uuid", app.requireActivatedUser(app.showHousesHandler))
	router.HandlerFunc(http.MethodPut, "/v1/houses/:uuid", app.requireActivatedUser(app.updateHouseHandler))

	// tenants
	router.HandlerFunc(http.MethodGet, "/v1/tenants", app.requireActivatedUser(app.listTenantsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tenants", app.requireActivatedUser(app.createTenantHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tenants/:uuid", app.requireActivatedUser(app.showTenantsHandler))
	router.HandlerFunc(http.MethodPut, "/v1/houses/:uuid", app.requireActivatedUser(app.updateTenantsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/tenants/:uuid", app.requireActivatedUser(app.removeTenant))

	//do magic send bulk imports back
	//router.HandlerFunc(http.MethodPost, "/v1/bulk/houses", app.magicbulkHousesHandler)

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))

}
