package main

import (
	"errors"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) listPaymentsHandler(c echo.Context) error {

	payments, err := app.store.GetAllPayments(c.Request().Context())
	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payments": payments})
}

func (app *application) showPaymentHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	payment, err := app.store.GetPaymentById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
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
		TenantId  uuid.UUID `json:"tenant_id"`
		Period    int32     `json:"period"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Renewed   bool      `json:"renewed"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), input.TenantId)

	if err != nil {
		if err == db.ErrRecordNotFound {
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

	args := db.UpdateTenantParams{
		FirstName:      tenant.FirstName,
		LastName:       tenant.LastName,
		HouseID:        tenant.HouseID,
		Phone:          tenant.Phone,
		PersonalIDType: tenant.PersonalIDType,
		PersonalID:     tenant.PersonalID,
		Active:         tenant.Active,
		Sos:            tenant.Sos,
		Eos:            tenant.Eos,
		TenantID:       tenant.TenantID,
	}

	err = app.store.UpdateTenant(c.Request().Context(), args)

	if err != nil {

		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	payment, err := app.store.GetUnrenewedByTenantId(c.Request().Context(), input.TenantId)

	if err != nil {
		if err != db.ErrRecordNotFound {
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "payment not found"})

	}

	payment.Renewed = true

	params := db.UpdatePaymentParams{
		PaymentID: payment.PaymentID,
		Period:    payment.Period,
		StartDate: payment.StartDate,
		EndDate:   payment.EndDate,
		Renewed:   payment.Renewed,
	}
	err = app.store.UpdatePayment(c.Request().Context(), params)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	pay := db.CreatePaymentParams{
		TenantID:  input.TenantId,
		Period:    input.Period,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		Renewed:   input.Renewed,
	}

	err = app.store.CreatePayment(c.Request().Context(), pay)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payment": pay})

}

func (app *application) updatPaymentHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	payment, err := app.store.GetPaymentById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "payment not found"})

		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}

	}

	var input struct {
		Period    *int32     `json:"period"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
		Renewed   *bool      `json:"renewed"`
	}

	if err := c.Bind(&input); err != nil {
		return err
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

	args := db.UpdatePaymentParams{
		PaymentID: payment.PaymentID,
		Period:    payment.Period,
		StartDate: payment.StartDate,
		EndDate:   payment.EndDate,
		Renewed:   payment.Renewed,
	}

	err = app.store.UpdatePayment(c.Request().Context(), args)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"payment": payment})

}
