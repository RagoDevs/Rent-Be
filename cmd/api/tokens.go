package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hopertz/rmgmt/internal/data"
	"github.com/Hopertz/rmgmt/internal/validator"
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

	v := validator.New()

	data.ValidateEmail(v, input.Email)

	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
	}

	var admin *data.Admin

	admin, err := app.models.Admins.GetByEmail(input.Email)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusUnauthorized, envelope{"error": "invalid email address or password"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	match, err := admin.Password.Matches(input.Password)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})

	}

	if !match {
		return c.JSON(http.StatusUnauthorized, envelope{"error": "invalid email address or password"})
	}

	token, err := app.models.Tokens.New(admin.AdminID, 3*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, envelope{"token": token.Plaintext})
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

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
	}
	// Try to retrieve the corresponding user record for the email address. If it can't
	// be found, return an error message to the client.
	admin, err := app.models.Admins.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}
	// Return an error message if the user is not activated.
	if !admin.Activated {
		v.AddError("email", "admin account must be activated")
		return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
	}
	// Otherwise, create a new password reset token with a 45-minute expiry time.
	token, err := app.models.Tokens.New(admin.AdminID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}
	// Email the user with their password reset token.
	app.background(func() {
		data := map[string]interface{}{
			"passwordResetToken": token.Plaintext,
		}
		// Since email addresses MAY be case sensitive, notice that we are sending this
		// email using the address stored in our database for the user --- not to the
		// input.Email address provided by the client in this request.
		err = app.mailer.Send(admin.Email, "token_password_reset.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	// Send a 202 Accepted response and confirmation message to the client.
	env := envelope{"message": "an email will be sent to you containing password reset instructions"}

	return c.JSON(http.StatusAccepted, env)
}

func (app *application) createActivationTokenHandler(c echo.Context) error {
	// Parse and validate the user's email address.
	var input struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
	}
	// Try to retrieve the corresponding user record for the email address. If it can't
	// be found, return an error message to the client.
	admin, err := app.models.Admins.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}
	// Return an error if the user has already been activated.
	if admin.Activated {
		v.AddError("email", "admin has already been activated")
		return c.JSON(http.StatusUnprocessableEntity, envelope{"errors": v.Errors})
	}
	// Otherwise, create a new activation token.
	token, err := app.models.Tokens.New(admin.AdminID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}
	// Email the user with their additional activation token.
	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
		}
		// Since email addresses MAY be case sensitive, notice that we are sending this
		// email using the address stored in our database for the user --- not to the
		// input.Email address provided by the client in this request.
		err = app.mailer.Send(admin.Email, "token_activation.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})
	// Send a 202 Accepted response and confirmation message to the client.
	env := envelope{"message": "an email will be sent to you containing activation instructions"}
	return c.JSON(http.StatusAccepted, env)
}
