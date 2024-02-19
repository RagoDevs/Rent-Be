package main

import (
	"errors"
	"net/http"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) listTenantsHandler(c echo.Context) error {

	tenants, err := app.store.GetTenants(c.Request().Context())

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenants": tenants})

}

func (app *application) showTenantsHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "tenant not found"})

		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenant": tenant})
}

func (app *application) createTenantHandler(c echo.Context) error {

	var input struct {
		FirstName      string    `json:"first_name"`
		LastName       string    `json:"last_name"`
		Phone          string    `json:"phone"`
		HouseId        uuid.UUID `json:"house_id"`
		PersonalIdType string    `json:"personal_id_type"`
		PersonalId     string    `json:"personal_id"`
		Active         bool      `json:"active"`
		Sos            time.Time `json:"sos"`
		Eos            time.Time `json:"eos"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	house, err := app.store.GetHouseById(c.Request().Context(), input.HouseId)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "house not found"})
	}
	if house.Occupied {
		return c.JSON(http.StatusConflict, map[string]interface{}{"error": "house is already occupied"})
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
		Eos:            input.Eos,
	}

	err = app.store.CreateTenant(c.Request().Context(), args)
	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	params := db.UpdateHouseByIdParams{
		Occupied: true,
		HouseID:  input.HouseId,
	}
	err = app.store.UpdateHouseById(c.Request().Context(), params)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})

	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"tenant": args})

}

func (app *application) updateTenantsHandler(c echo.Context) error {

	id, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), id)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "tenant not found"})

		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
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
		Eos            *time.Time `json:"eos"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

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

	if input.Eos != nil {
		tenant.Eos = *input.Eos
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
		TenantID:       tenant.TenantID,
	}
	err = app.store.UpdateTenant(c.Request().Context(), arg)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenant": tenant})

}

func (app *application) removeTenant(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return err
	}

	tenant, err := app.store.GetTenantById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, db.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "tenant not found"})
		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	tenant.Active = false

	err = app.store.UpdateTenant(c.Request().Context(), db.UpdateTenantParams{
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
	})

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	err = app.store.UpdateHouseById(c.Request().Context(), db.UpdateHouseByIdParams{
		Occupied: false,
		HouseID:  tenant.HouseID,
	})

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenant": tenant})

}
