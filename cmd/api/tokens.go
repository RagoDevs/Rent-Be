package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rent/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) createAuthenticationTokenHandler(c echo.Context) error {

	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	admin, err := app.store.GetAdminByPhone(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by email", "error", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid email number or password"})
		default:
			slog.Error("error fetching admin by phone number", "error", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	if !admin.Activated {
		return c.JSON(http.StatusBadRequest, envelope{"error": "admin not activated"})
	}

	pwd := db.Password{
		Hash:      admin.PasswordHash,
		Plaintext: input.Password,
	}

	match, err := db.PasswordMatches(pwd)

	if err != nil {
		slog.Error("error matching password", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})

	}

	if !match {
		slog.Error("error matching password", "error", err)
		return c.JSON(http.StatusUnauthorized, envelope{"error": "invalid phone number or password"})
	}

	expiry := time.Now().Add(3 * 24 * time.Hour)
	token, err := app.store.NewToken(admin.ID, expiry, db.ScopeAuthentication)
	if err != nil {
		slog.Error("error generating new token", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, envelope{"token": token.Plaintext, "expiry": expiry})
}

func (app *application) createPasswordResetTokenHandler(c echo.Context) error {

	var input struct {
		Phone string `json:"phone" validate:"required,len=10"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	admin, err := app.store.GetAdminByPhone(c.Request().Context(), input.Phone)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by phone number", "error", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid phone number or password"})
		default:
			slog.Error("error fetching admin by phone number", "error", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	if !admin.Activated {
		return c.JSON(http.StatusForbidden, envelope{"errors": "account not activated"})
	}

	expiry := time.Now().Add(45 * time.Minute)

	token, err := app.store.NewToken(admin.ID, expiry, db.ScopePasswordReset)
	if err != nil {
		slog.Error("error generating new token", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	msg := fmt.Sprintf("Your password reset token is: %s", token.Plaintext)

	app.background(func() {

		// send mail here

		_ = msg
	})

	return c.JSON(http.StatusAccepted, nil)
}

func (app *application) createActivationTokenHandler(c echo.Context) error {

	var input struct {
		Phone string `json:"phone" validate:"required,len=10"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	admin, err := app.store.GetAdminByPhone(c.Request().Context(), input.Phone)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by phone number", "error", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid phone number or password"})
		default:
			slog.Error("error fetching admin by phone number", "error", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	if admin.Activated {
		return c.JSON(http.StatusForbidden, envelope{"errors": "account already activated"})
	}

	expiry := time.Now().Add(3 * 24 * time.Hour)

	token, err := app.store.NewToken(admin.ID, expiry, db.ScopeActivation)
	if err != nil {
		slog.Error("error generating new token", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	msg := fmt.Sprintf("Your account activation token is: %s", token.Plaintext)

	app.background(func() {

		_ = msg
		//send mail here
	})

	return c.JSON(http.StatusAccepted, nil)
}
