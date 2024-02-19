package main

import (
	"crypto/sha256"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) registerUserHandler(c echo.Context) error {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	// validate email above?

	pwd, err := db.SetPassword(input.Password)

	if err != nil {
		return err
	}

	a := db.InsertAdminParams{
		Email:        input.Email,
		PasswordHash: pwd.Hash,
		Activated:    false,
	}

	res, err := app.store.InsertAdmin(c.Request().Context(), a)
	if err != nil {
		switch {

		case errors.Is(err, db.ErrDuplicateEmail):
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})

		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

	}

	token, err := app.store.NewToken(res.AdminID, 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		//logerror above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"adminID":         res.AdminID,
		}

		err = app.mailer.Send(input.Email, "admin_welcome.tmpl", data)
		if err != nil {
			slog.Error("error sending email", err)
		}
	})

	return c.JSON(http.StatusCreated, map[string]interface{}{"message": "admin created successfully"})
}

func (app *application) activateUserHandler(c echo.Context) error {

	var input struct {
		TokenPlaintext string `json:"token"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if Isvalid, err := db.IsValidTokenPlaintext(input.TokenPlaintext); !Isvalid {
		return c.JSON(http.StatusUnprocessableEntity, err)
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	args := db.GetForTokenAdminParams{
		Scope:  db.ScopeActivation,
		Hash:   tokenHash[:],
		Expiry: time.Now(),
	}

	admin, err := app.store.GetForTokenAdmin(c.Request().Context(), args)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
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
		case errors.Is(err, db.ErrEditConflict):
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "unable to complete request due to an edit conflict"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

	}

	a := db.DeleteAllTokenParams{
		Scope:   db.ScopeActivation,
		AdminID: admin.AdminID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), a)

	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Your account has been activated successfully"})
}

// Verify the password reset token and set a new password for the admin.
func (app *application) updateUserPasswordHandler(c echo.Context) error {
	// Parse and validate the admins's new password and password reset token.
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	tokenHash := sha256.Sum256([]byte(input.TokenPlaintext))

	args := db.GetForTokenAdminParams{
		Scope:  db.ScopePasswordReset,
		Hash:   tokenHash[:],
		Expiry: time.Now(),
	}

	admin, err := app.store.GetForTokenAdmin(c.Request().Context(), args)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):

			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}
	// Set the new password for the admin.

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
		case errors.Is(err, db.ErrEditConflict):
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "unable to complete request due to an edit conflict"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	d := db.DeleteAllTokenParams{
		Scope:   db.ScopePasswordReset,
		AdminID: admin.AdminID,
	}
	err = app.store.DeleteAllToken(c.Request().Context(), d)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Your password has been updated successfully"})
}
