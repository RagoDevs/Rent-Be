package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/labstack/echo/v4"
)

func (app *application) createAuthenticationTokenHandler(c echo.Context) error {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid email address or password"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	pwd := db.Password{
		Hash:      admin.PasswordHash,
		Plaintext: input.Password,
	}

	match, err := db.PasswordMatches(pwd)

	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})

	}

	if !match {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "invalid email address or password"})
	}

	token, err := app.store.NewToken(admin.AdminID, 3*24*time.Hour, db.ScopeAuthentication)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"token": token.Plaintext})
}

// Generate a password reset token and send it to the user's email address.
func (app *application) createPasswordResetTokenHandler(c echo.Context) error {
	// Parse and validate the user's email address.
	var input struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	if !admin.Activated {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})
	}

	token, err := app.store.NewToken(admin.AdminID, 45*time.Minute, db.ScopePasswordReset)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	app.background(func() {
		data := map[string]interface{}{
			"passwordResetToken": token.Plaintext,
		}
		err = app.mailer.Send(admin.Email, "token_password_reset.tmpl", data)
		if err != nil {
			slog.Error("error sending ", err)
		}
	})

	env := map[string]interface{}{"message": "an email will be sent to you containing password reset instructions"}

	return c.JSON(http.StatusAccepted, env)
}

func (app *application) createActivationTokenHandler(c echo.Context) error {

	var input struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	admin, err := app.store.GetAdminByEmail(c.Request().Context(), input.Email)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors":""})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	if admin.Activated {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": ""})
	}

	token, err := app.store.NewToken(admin.AdminID, 3*24*time.Hour, db.ScopeActivation)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
		}

		err = app.mailer.Send(admin.Email, "token_activation.tmpl", data)
		if err != nil {
			slog.Error("error sending email", err)
		}
	})
	// Send a 202 Accepted response and confirmation message to the client.
	env := map[string]interface{}{"message": "an email will be sent to you containing activation instructions"}
	return c.JSON(http.StatusAccepted, env)
}
