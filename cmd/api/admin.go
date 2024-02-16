package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hopertz/rmgmt/internal/data"
	"github.com/Hopertz/rmgmt/internal/validator"
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

	admin := &data.Admin{
		Email:     input.Email,
		Activated: false,
	}

	err := admin.Password.Set(input.Password)
	if err != nil {
		return err
	}

	v := validator.New()

	if data.ValidateUser(v, admin); !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})
	}

	err = app.models.Admins.Insert(admin)
	if err != nil {
		switch {

		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "an admin with this email address already exists")
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})

		default:
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
		
	}

	token, err := app.models.Tokens.New(admin.AdminID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		//logerror above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	app.background(func() {
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"adminID":         admin.AdminID,
		}

		err = app.mailer.Send(admin.Email, "admin_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
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

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})
	}

	admin, err := app.models.Admins.GetForToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
		
	}

	admin.Activated = true

	err = app.models.Admins.Update(admin)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "unable to complete request due to an edit conflict"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
		
	}

	err = app.models.Tokens.DeleteAllForAdmin(data.ScopeActivation, admin.AdminID)
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

	v := validator.New()
	data.ValidatePasswordPlaintext(v, input.Password)
	data.ValidateTokenPlaintext(v, input.TokenPlaintext)
	if !v.Valid() {
		return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})
	}
	// Retrieve the details of the admin associated with the password reset token,
	// returning an error message if no matching record was found.
	// returning an error message if no matching record was found.
	admin, err := app.models.Admins.GetForToken(data.ScopePasswordReset, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "Invalid or expired password reset token")
			return c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{"errors": v.Errors})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}
	// Set the new password for the admin.
	err = admin.Password.Set(input.Password)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}
	// Save the updated user record in our database, checking for any edit conflicts as
	// normal.
	err = app.models.Admins.Update(admin)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "unable to complete request due to an edit conflict"})
		default:
			//log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}
	// If everything was successful, then delete all password reset tokens for the admin.
	err = app.models.Tokens.DeleteAllForAdmin(data.ScopePasswordReset, admin.AdminID)
	if err != nil {
		//log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Your password has been updated successfully"})
}
