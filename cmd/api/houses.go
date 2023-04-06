package main

import (
	"errors"
	"net/http"

	"hmgt.hopertz.me/internal/data"
)

func (app *application) listHousesHandler(w http.ResponseWriter, r *http.Request) {
	houses, err := app.models.Houses.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"houses": houses}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showHousesHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	house, err := app.models.Houses.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"house": house}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createHouseHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Location  string `json:"location"`
		Block     string `json:"block"`
		Partition string `json:"partition"`
		Occupied  bool   `json:"occupied"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	house := &data.House{
		Location:  input.Location,
		Block:     input.Block,
		Partition: input.Partition,
		Occupied:  input.Occupied,
	}

	err = app.models.Houses.Insert(house)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"house": house}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}



func (app *application) updateHouseHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	house, err := app.models.Houses.Get(uuid)

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
		
		Occupied  bool   `json:"occupied"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	

	house.Occupied = input.Occupied
	

	err = app.models.Houses.Update(house)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"house": house}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
