package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"time"

	db "github.com/Hopertz/rent/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) listPaymentsHandler(c echo.Context) error {

	payments, err := app.store.GetAllPayments(c.Request().Context())
	if err != nil {
		slog.Error("error fetching payments", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, payments)
}

func (app *application) showPaymentHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	payment, err := app.store.GetDetailedPaymentById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusNotFound, envelope{"error": "payment not found"})

		default:
			slog.Error("error fetching payment by id", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, payment)
}

func (app *application) createPaymentHandler(c echo.Context) error {

	admin, ok := c.Get("admin").(db.GetHashTokenForAdminRow)

	if !ok {
		return c.JSON(http.StatusBadRequest, envelope{"error": "you are not authorized to perform this action"})
	}

	var input struct {
		TenantId  uuid.UUID `json:"tenant_id" validate:"required"`
		Amount    int32     `json:"amount" validate:"required"`
		StartDate time.Time `json:"start_date" validate:"required"`
		EndDate   time.Time `json:"end_date" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), input.TenantId)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, envelope{"error": "tenant not found"})
		}

		slog.Error("error fetching tenant by id in create payments", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})

	}

	err = app.store.CreatePayment(c.Request().Context(), db.CreatePaymentParams{
		TenantID:  tenant.HouseID,
		Amount:    input.Amount,
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
		CreatedBy: admin.ID,
	})

	if err != nil {
		slog.Error("error creating payment", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (app *application) updatePaymentHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid payment id"})
	}

	payment, err := app.store.GetPaymentById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching payment by id on update payment", "error", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "payment not found"})

		default:
			slog.Error("error fetching payment by id on update payment", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	var input struct {
		Amount    *int32     `json:"amount"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	if input.Amount != nil {
		payment.Amount = *input.Amount
	}

	if input.StartDate != nil {
		payment.StartDate = *input.StartDate
	}

	if input.EndDate != nil {
		payment.EndDate = *input.EndDate
	}

	args := db.UpdatePaymentParams{
		ID:        payment.ID,
		Amount:    payment.Amount,
		StartDate: payment.StartDate,
		EndDate:   payment.EndDate,
		Version:   payment.Version,
	}

	err = app.store.UpdatePayment(c.Request().Context(), args)

	if err != nil {
		slog.Error("error updating payment", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)

}

func (app *application) deletePaymentHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid payment id"})
	}

	payment, err := app.store.GetPaymentById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching payment by id on delete payment", "error", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "payment not found"})

		default:
			slog.Error("error fetching payment by id on delete payment", "error", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	err = app.store.DeletePayment(c.Request().Context(), payment.ID)

	if err != nil {
		slog.Error("error deleting payment", "error", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}
