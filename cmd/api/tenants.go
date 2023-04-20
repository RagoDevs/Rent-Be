package main

import (
	"errors"
	"net/http"
	"time"

	"hmgt.hopertz.me/internal/data"
)

func (app *application) listTenantsHandler(w http.ResponseWriter, r *http.Request) {
	tenants, err := app.models.Tenants.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenants": tenants}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showTenantsHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": tenant}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTenantHandler(w http.ResponseWriter, r *http.Request) {
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

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
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
		app.serverErrorResponse(w, r, err)
	}

	err = app.models.Houses.Update(tenant.HouseId, true)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": tenant}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateTenantsHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
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

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
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
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tenant": tenant}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) removeTenant(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tenant, err := app.models.Tenants.Get(uuid)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	tenant.Active = false

	err = app.models.Tenants.Update(tenant)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.models.Houses.Update(tenant.HouseId, false)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
