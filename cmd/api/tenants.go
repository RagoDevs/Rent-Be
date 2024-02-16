package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Hopertz/rmgmt/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) listTenantsHandler(c echo.Context) error {

	tenants, err := app.models.Tenants.GetAll()

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenants": tenants})

}

func (app *application) showTenantsHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
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
		HouseId        string    `json:"house_id"`
		PersonalIdType string    `json:"personal_id_type"`
		PersonalId     string    `json:"personal_id"`
		Active         bool      `json:"active"`
		Sos            time.Time `json:"sos"`
		Eos            time.Time `json:"eos"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	house, err := app.models.Houses.Get(input.HouseId)

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "house not found"})
	}
	if house.Occupied {
		return c.JSON(http.StatusConflict, map[string]interface{}{"error": "house is already occupied"})
	}

	tenant := &data.Tenant{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Phone:          input.Phone,
		HouseId:        input.HouseId,
		PersonalIdType: input.PersonalIdType,
		PersonalId:     input.PersonalId,
		Active:         input.Active,
		Sos:            input.Sos,
		Eos:            input.Eos,
	}

	err = app.models.Tenants.Insert(tenant)

	if err != nil {
		if err == data.ErrDuplicatePhoneNumber {
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "duplicate phone number"})
		}

		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	err = app.models.Houses.Update(tenant.HouseId, true)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})

	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"tenant": tenant})

}

func (app *application) updateTenantsHandler(c echo.Context) error {

	uuid := c.Param("uuid")

	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
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
		HouseId        *string    `json:"house_id"`
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
		tenant.HouseId = *input.HouseId
	}

	if input.PersonalIdType != nil {
		tenant.PersonalIdType = *input.PersonalIdType
	}

	if input.PersonalId != nil {
		tenant.PersonalId = *input.PersonalId
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

	err = app.models.Tenants.Update(tenant)

	if err != nil {
		if err == data.ErrDuplicatePhoneNumber {
			return c.JSON(http.StatusConflict, map[string]interface{}{"error": "duplicate phone number"})
		}

		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenant": tenant})

}

func (app *application) removeTenant(c echo.Context) error {
	uuid := c.Param("uuid")

	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "tenant not found"})
		default:
			// log error above
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
		}
	}

	tenant.Active = false

	err = app.models.Tenants.Update(tenant)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	err = app.models.Houses.Update(tenant.HouseId, false)

	if err != nil {
		// log error above
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"tenant": tenant})

}
