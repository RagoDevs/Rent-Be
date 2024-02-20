package main

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) registerUserHandler(c echo.Context) error {

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

	pwd, err := db.SetPassword(input.Password)

	if err != nil {
		slog.Error("error generating hash password", err)
		return err
	}

	args := db.InsertAdminParams{
		Email:        input.Email,
		PasswordHash: pwd.Hash,
		Activated:    false,
	}

	a, err := app.store.InsertAdmin(c.Request().Context(), args)

	if err != nil {
		switch {

		case err.Error() == db.DuplicateEmail:
			return c.JSON(http.StatusBadRequest, envelope{"error": "email address is already in use"})

		default:
			slog.Error("error creating admin", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	token, err := app.store.NewToken(a.AdminID, 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		slog.Error("error generating new token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"id":              a.AdminID,
		}

		err = app.mailer.Send(input.Email, "admin_welcome.tmpl", data)
		if err != nil {
			slog.Error("error sending email", err)
		}
	})

	return c.JSON(http.StatusCreated, nil)
}

func (app *application) activateUserHandler(c echo.Context) error {

	var input struct {
		TokenPlaintext string `json:"token" validate:"required,len=26"`
	}

	if err := c.Bind(&input); err != nil {
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

		AdminID:      admin.AdminID,
		Email:        admin.Email,
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
		Scope:   db.ScopeActivation,
		AdminID: admin.AdminID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), a)

	if err != nil {
		slog.Error("error deleting token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}

func (app *application) updateUserPasswordHandler(c echo.Context) error {

	var input struct {
		Password string `json:"password" validate:"required,min=8"`
		TokenPlaintext string `json:"token" validate:"required,len=26"`
	}

	if err := c.Bind(&input); err != nil {
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
		Email:        admin.Email,
		PasswordHash: pwd.Hash,
		Activated:    true,
		AdminID:      admin.AdminID,
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
		Scope:   db.ScopePasswordReset,
		AdminID: admin.AdminID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), d)
	
	if err != nil {
		slog.Error("error deleting token", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}
