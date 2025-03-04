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

func (app *application) listTenantsHandler(c echo.Context) error {

	tenants, err := app.store.GetTenants(c.Request().Context())

	if err != nil {
		slog.Error("error fetching tenants", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, tenants)

}

func (app *application) showTenantsHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching tenant by id", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "tenant not found"})

		default:
			slog.Error("error fetching tenant by id", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, tenant)
}

func (app *application) createTenantHandler(c echo.Context) error {

	var input struct {
		FirstName      string    `json:"first_name" validate:"required"`
		LastName       string    `json:"last_name" validate:"required"`
		Phone          string    `json:"phone" validate:"required,len=10"`
		HouseId        uuid.UUID `json:"house_id" validate:"required"`
		PersonalIdType string    `json:"personal_id_type" validate:"required"`
		PersonalId     string    `json:"personal_id" validate:"required"`
		Active         bool      `json:"active"`
		Sos            time.Time `json:"sos" validate:"required"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	house, err := app.store.GetHouseById(c.Request().Context(), input.HouseId)

	if err != nil {

		switch {

		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching house by id for create tenant", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "house not found"})

		default:
			slog.Error("error fetching house by id for create tenant", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	if house.Occupied {
		return c.JSON(http.StatusBadRequest, envelope{"error": "house is already occupied"})
	}

	args := db.CreateTenantParams{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		HouseID:        input.HouseId,
		Phone:          input.Phone,
		PersonalIDType: input.PersonalIdType,
		PersonalID:     input.PersonalId,
		Active:         input.Active,
		Sos:            input.Sos,
	}

	err = app.store.TxnCreateTenant(c.Request().Context(), args)

	if err != nil {
		slog.Error("error creating tenant", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, nil)

}

func (app *application) updateTenantsHandler(c echo.Context) error {

	id, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching tenant by id for update", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "tenant not found"})

		default:
			slog.Error("error fetching tenant by id for update", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	var input struct {
		FirstName      *string    `json:"first_name"`
		LastName       *string    `json:"last_name"`
		Phone          *string    `json:"phone"`
		HouseId        *uuid.UUID `json:"house_id"`
		PersonalIdType *string    `json:"personal_id_type"`
		PersonalId     *string    `json:"personal_id"`
		Active         *bool      `json:"active"`
		Sos            *time.Time `json:"sos"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	prev_house_id := tenant.HouseID

	if input.FirstName != nil {
		tenant.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		tenant.LastName = *input.LastName
	}

	if input.Phone != nil {
		tenant.Phone = *input.Phone
	}

	if input.HouseId != nil {
		tenant.HouseID = *input.HouseId
	}

	if input.PersonalIdType != nil {
		tenant.PersonalIDType = *input.PersonalIdType
	}

	if input.PersonalId != nil {
		tenant.PersonalID = *input.PersonalId
	}

	if input.Active != nil {
		tenant.Active = *input.Active
	}

	if input.Sos != nil {
		tenant.Sos = *input.Sos
	}

	arg := db.UpdateTenantParams{

		FirstName:      tenant.FirstName,
		LastName:       tenant.LastName,
		HouseID:        tenant.HouseID,
		Phone:          tenant.Phone,
		PersonalIDType: tenant.PersonalIDType,
		PersonalID:     tenant.PersonalID,
		Active:         tenant.Active,
		Sos:            tenant.Sos,
		Eos:            tenant.Eos,
		ID:             tenant.ID,
		Version:        tenant.Version,
	}

	err = app.store.TxnUpdateTenantHouse(c.Request().Context(), arg, prev_house_id)

	if err != nil {
		slog.Error("error updating tenant and house", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)

}

func (app *application) removeTenant(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching tenant by id for update", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "tenant not found"})
		default:
			slog.Error("error fetching tenant by id for update", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	tenant.Active = false

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
		ID:             tenant.ID,
		Version:        tenant.Version,
	}

	err = app.store.TxnRemoveTenantHouse(c.Request().Context(), args)

	if err != nil {
		slog.Error("failed deactiving tenant & disabling house", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)

}
