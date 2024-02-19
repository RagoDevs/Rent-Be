package main

import (
	"crypto/sha256"
	"errors"
	"strings"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/Hopertz/rmgmt/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *application) authenticate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		authorizationHeader := c.Request().Header.Get("Authorization")

		if authorizationHeader == "" {

			c.Set("user", db.Admin{})

			return next(c)
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

			return c.JSON(401, "invalid or missing authentication token")

		}

		token := headerParts[1]
		v := validator.New()

		if db.ValidateTokenPlaintext(v, token); !v.Valid() {
			return c.JSON(401, "invalid or missing authentication token")
		}

		tokenHash := sha256.Sum256([]byte(token))
		args := db.GetForTokenAdminParams{
			Scope:  db.ScopeAuthentication,
			Hash:   tokenHash[:],
			Expiry: time.Now(),
		}

		admin, err := app.store.GetForTokenAdmin(c.Request().Context(), args)
		if err != nil {
			switch {
			case errors.Is(err, db.ErrRecordNotFound):
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
		user, ok := c.Get("admin").(*db.Admin)

		if !ok {
			return c.JSON(401, "unauthorized you must be authencticated ")
		}

		// if user.IsAnonymous() {
		// 	return c.JSON(401, "unauthorized no anonymous user allowed")
		// }

		if !user.Activated {
			return c.JSON(403, "forbidden user account not activated")
		}

		return next(c)
	}
}
