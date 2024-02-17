package main

import (
	"errors"
	"strings"

	"github.com/Hopertz/rmgmt/db/data"
	"github.com/Hopertz/rmgmt/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *application) authenticate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		authorizationHeader := c.Request().Header.Get("Authorization")

		if authorizationHeader == "" {

			c.Set("user", data.AnonymousAdmin)

			return next(c)
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

			return c.JSON(401, "invalid or missing authentication token")

		}

		token := headerParts[1]
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			return c.JSON(401, "invalid or missing authentication token")
		}

		admin, err := app.models.Admins.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				return c.JSON(401, "invalid or missing authentication token")
			default:
				return c.JSON(500, "pkg server error")
			}
		}

		c.Set("admin", admin)

		return next(c)

	}

}

func (app *application) requireAuthenticatedUser(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		user, ok := c.Get("admin").(*data.Admin)

		if !ok {
			return c.JSON(401, "unauthorized you must be authencticated ")
		}

		if user.IsAnonymous() {
			return c.JSON(401, "unauthorized no anonymous user allowed")
		}

		if !user.Activated {
			return c.JSON(403, "forbidden user account not activated")
		}

		return next(c)
	}
}
