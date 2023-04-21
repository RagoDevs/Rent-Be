package main

import (
	"errors"
	"net/http"
	"time"

	"hmgt.hopertz.me/internal/data"
)

func (app *application) listPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	payments, err := app.models.Payments.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payments": payments}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showPaymentHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	payment, err := app.models.Payments.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payment": payment}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TenantId  string    `json:"tenant_id"`
		Period    int       `json:"period"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
		Renewed   bool      `json:"renewed"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	_, err = app.models.Tenants.Get(input.TenantId)

	if err != nil {
		if err == data.ErrRecordNotFound {
			app.badRequestResponse(w, r, errors.New("tenant cannot be found"))
			return
		}

	}

	payment, err := app.models.Payments.GetUnrenewed(input.TenantId)

	if err != nil {
		if err != data.ErrRecordNotFound {
			app.serverErrorResponse(w, r, err)
			return
		}

	}

	if payment != nil {
		payment.Renewed = true
		err = app.models.Payments.Update(*payment)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payment": payment}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updatPaymentHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	payment, err := app.models.Payments.Get(uuid)

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
		TenantId  *string    `json:"tenant_id"`
		Period    *int       `json:"period"`
		StartDate *time.Time `json:"start_date"`
		EndDate   *time.Time `json:"end_date"`
		Renewed   *bool      `json:"renewed"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"payment": payment}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}


