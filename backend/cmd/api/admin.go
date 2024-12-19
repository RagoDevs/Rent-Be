package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rent/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) registerAdminHandler(c echo.Context) error {

	var input struct {
		Phone    string `json:"phone" validate:"required,len=10"`
		Password string `json:"password" validate:"required,min=8"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if input.Phone != app.config.phone {
		return c.JSON(http.StatusUnauthorized, envelope{"error": "phone not allowed"})
	}

	pwd, err := db.SetPassword(input.Password)

	if err != nil {
		slog.Error("error generating hash password", err)
		return err
	}

	args := db.CreateAdminParams{
		Phone:        input.Phone,
		PasswordHash: pwd.Hash,
		Activated:    false,
	}

	a, err := app.store.CreateAdmin(c.Request().Context(), args)

	if err != nil {
		switch {

		case err.Error() == db.DuplicatePhone:
			return c.JSON(http.StatusBadRequest, envelope{"error": "phone number is already in use"})

		default:
			slog.Error("error creating admin", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	token, err := app.store.NewToken(a.ID, 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		slog.Error("error generating new token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	msg := fmt.Sprintf("Welcome, your activation token is %s", token.Plaintext)

	app.background(func() {

		err = app.beem.Send(msg, input.Phone)

		if err != nil {
			slog.Error("error sending ", "err", err)
		}
	})

	return c.JSON(http.StatusCreated, nil)
}

func (app *application) activateAdminHandler(c echo.Context) error {

	var input struct {
		TokenPlaintext string `json:"token" validate:"required,len=26"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	args := db.GetHashTokenForAdminParams{
		Scope:  db.ScopeActivation,
		Hash:   tokenHash[:],
		Expiry: time.Now(),
	}

	admin, err := app.store.GetHashTokenForAdmin(c.Request().Context(), args)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching token user admin", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "invalid token"})
		default:
			slog.Error("error fetching token user admin", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	admin.Activated = true

	param := db.UpdateAdminParams{

		ID:           admin.ID,
		Phone:        admin.Phone,
		Activated:    true,
		PasswordHash: admin.PasswordHash,
		Version:      admin.Version,
	}
	_, err = app.store.UpdateAdmin(c.Request().Context(), param)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error conflict updating admin ", err)
			return c.JSON(http.StatusConflict, envelope{"error": "unable to complete request due to an edit conflict"})
		default:
			slog.Error("error updating admin ", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	a := db.DeleteAllTokenParams{
		Scope: db.ScopeActivation,
		ID:    admin.ID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), a)

	if err != nil {
		slog.Error("error deleting token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}

func (app *application) updateAdminPasswordHandler(c echo.Context) error {

	var input struct {
		Password       string `json:"password" validate:"required,min=8"`
		TokenPlaintext string `json:"token" validate:"required,len=26"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	args := db.GetHashTokenForAdminParams{
		Scope:  db.ScopePasswordReset,
		Hash:   tokenHash[:],
		Expiry: time.Now(),
	}

	admin, err := app.store.GetHashTokenForAdmin(c.Request().Context(), args)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching token user admin", err)
			return c.JSON(http.StatusNotFound, envelope{"errors": "invalid token"})
		default:
			slog.Error("error fetching token user admin", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	pwd, err := db.SetPassword(input.Password)

	if err != nil {
		return err
	}

	_, err = app.store.UpdateAdmin(c.Request().Context(), db.UpdateAdminParams{
		Phone:        admin.Phone,
		PasswordHash: pwd.Hash,
		Activated:    true,
		ID:           admin.ID,
		Version:      admin.Version,
	})

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusConflict, envelope{"error": "unable to complete request due to an edit conflict"})
		default:
			slog.Error("error updating admin ", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	d := db.DeleteAllTokenParams{
		Scope: db.ScopePasswordReset,
		ID:    admin.ID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), d)

	if err != nil {
		slog.Error("error deleting token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}
