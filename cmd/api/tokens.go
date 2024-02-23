package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
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

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid email address or password"})
		default:
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	pwd := db.Password{
		Hash:      admin.PasswordHash,
		Plaintext: input.Password,
	}

	match, err := db.PasswordMatches(pwd)

	if err != nil {
		slog.Error("error matching password", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})

	}

	if !match {
		slog.Error("error matching password", err)
		return c.JSON(http.StatusUnauthorized, envelope{"error": "invalid email address or password"})
	}

	token, err := app.store.NewToken(admin.ID, 3*24*time.Hour, db.ScopeAuthentication)
	if err != nil {
		slog.Error("error generating new token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, envelope{"token": token.Plaintext})
}

func (app *application) createPasswordResetTokenHandler(c echo.Context) error {

	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid email address or password"})
		default:
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	if !admin.Activated {
		return c.JSON(http.StatusForbidden, envelope{"errors": "account not activated"})
	}

	token, err := app.store.NewToken(admin.ID, 45*time.Minute, db.ScopePasswordReset)
	if err != nil {
		slog.Error("error generating new token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	app.background(func() {
		data := envelope{
			"passwordResetToken": token.Plaintext,
		}
		err = app.mailer.Send(admin.Email, "token_password_reset.tmpl", data, "Password Reset")
		if err != nil {
			slog.Error("error sending ", err)
		}
	})

	return c.JSON(http.StatusAccepted, nil)
}

func (app *application) createActivationTokenHandler(c echo.Context) error {

	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid email address or password"})
		default:
			slog.Error("error fetching admin by email", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	if admin.Activated {
		return c.JSON(http.StatusForbidden, envelope{"errors": "account already activated"})
	}

	token, err := app.store.NewToken(admin.ID, 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		slog.Error("error generating new token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	app.background(func() {
		data := envelope{
			"activationToken": token.Plaintext,
		}

		err = app.mailer.Send(admin.Email, "token_activation.tmpl", data, "Account Activation")
		if err != nil {
			slog.Error("error sending email", err)
		}
	})

	return c.JSON(http.StatusAccepted, nil)
}
