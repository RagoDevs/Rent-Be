package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"expvar"
	"log/slog"
	"net/http"
	"strings"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) authenticate(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		authorizationHeader := c.Request().Header.Get("Authorization")

		if authorizationHeader == "" {
			return c.JSON(http.StatusBadRequest, envelope{"error": "missing authentication token"})
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {

			return c.JSON(http.StatusBadRequest, envelope{"error": "invalid token"})

		}

		token := headerParts[1]

		if Isvalid, err := db.IsValidTokenPlaintext(token); !Isvalid {
			return c.JSON(http.StatusBadRequest, envelope{"error": err})
		}

		tokenHash := sha256.Sum256([]byte(token))

		args := db.GetHashTokenForAdminParams{
			Scope:  db.ScopeAuthentication,
			Hash:   tokenHash[:],
			Expiry: time.Now(),
		}

		admin, err := app.store.GetHashTokenForAdmin(c.Request().Context(), args)

		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				slog.Error("error", err)
				return c.JSON(http.StatusNotFound, envelope{"error": "invalid token"})
			default:
				slog.Error("error", err)
				return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
			}
		}

		c.Set("admin", admin)

		return next(c)

	}

}

func (app *application) requireAuthenticatedAdmin(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {

		admin, ok := c.Get("admin").(db.GetHashTokenForAdminRow)

		if !ok {
			return c.JSON(http.StatusBadRequest, envelope{"error": "missing admin from context"})
		}

		if !admin.Activated {
			return c.JSON(http.StatusForbidden, envelope{"error": "admin account not activated"})
		}

		return next(c)
	}
}

func (app *application) metrics(next echo.HandlerFunc) echo.HandlerFunc {

	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	return func(c echo.Context) error {

		totalRequestsReceived.Add(1)
		start := time.Now()

		err := next(c)

		totalResponsesSent.Add(1)
		totalProcessingTimeMicroseconds.Add(time.Since(start).Microseconds())

		return err
	}
}
