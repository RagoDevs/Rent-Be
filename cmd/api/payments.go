package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hopertz/rmgmt/db/data"
	"github.com/labstack/echo/v4"
)

func (app *application) listPaymentsHandler(c echo.Context) error {
	payments, err := app.models.Payments.GetAll()

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payments": payments})
}

func (app *application) showPaymentHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	payment, err := app.models.Payments.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "payment not found"})

		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payment": payment})
}

func (app *application) createPaymentHandler(c echo.Context) error {
	var input struct {
		TenantId  string    `json:"tenant_id"`
		Period    int       `json:"period"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Renewed   bool      `json:"renewed"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	tenant, err := app.models.Tenants.Get(input.TenantId)

	if err != nil {
		if err == data.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "tenant not found"})
		}

	}

	// IsEqual := tenant.Eos.Equal(input.EndDate)

	// IsBefore := tenant.Eos.Before(input.EndDate)

	// if !IsEqual || !IsBefore {

	// 	app.badRequestResponse(w, r, errors.New("end of stay should be less or equal to end date of payment"))

	// 	return
	// }

	tenant.Eos = input.EndDate

	err = app.models.Tenants.Update(tenant)

	if err != nil {

		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	payment, err := app.models.Payments.GetUnrenewed(input.TenantId)

	if err != nil {
		if err != data.ErrRecordNotFound {
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

	}

	if payment != nil {
		payment.Renewed = true
		err = app.models.Payments.Update(*payment)
		if err != nil {
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	payment = &data.Payment{

		TenantId:  input.TenantId,
		Period:    input.Period,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Renewed:   input.Renewed,
	}

	err = app.models.Payments.Insert(payment)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payment": payment})

}

func (app *application) updatPaymentHandler(c echo.Context) error {
	uuid := c.Param("uuid")

	payment, err := app.models.Payments.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "payment not found"})

		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

	}

	var input struct {
		TenantId  *string    `json:"tenant_id"`
		Period    *int       `json:"period"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
		Renewed   *bool      `json:"renewed"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if input.TenantId != nil {
		payment.TenantId = *input.TenantId
	}

	if input.Period != nil {
		payment.Period = *input.Period
	}

	if input.StartDate != nil {
		payment.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		payment.EndDate = *input.EndDate
	}

	if input.Renewed != nil {
		payment.Renewed = *input.Renewed
	}

	err = app.models.Payments.Update(*payment)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payment": payment})

}
