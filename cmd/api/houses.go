package main

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	db "github.com/Hopertz/rent/db/sqlc"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type FormattedHouseData struct {
	HouseID   uuid.UUID `json:"house_id"`
	Location  string    `json:"location"`
	Block     string    `json:"block"`
	Partition int16     `json:"partition"`
	Price     int32     `json:"price"`
	Occupied  bool      `json:"occupied"`
	Name      string    `json:"name"`
	TenantID  string    `json:"tenant_id"`
}

func (app *application) listHousesHandler(c echo.Context) error {

	houses, err := app.store.GetHouses(c.Request().Context())

	if err != nil {
		slog.Error("error fetching houses", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, houses)

}

func (app *application) showHouseHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid house id"})
	}

	house, err := app.store.GetHouseById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching house by id", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "house not found"})
		default:
			slog.Error("error fetching house by id", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	formattedHouse := FormattedHouseData{
		HouseID:   house.HouseID,
		Location:  house.Location,
		Block:     house.Block,
		Partition: house.Partition,
		Price:     house.Price,
		Occupied:  house.Occupied,
		Name:      "",
		TenantID:  "",
	}

	if house.Name.Valid {
		formattedHouse.Name = house.Name.String
	}

	if house.TenantID.Valid {
		formattedHouse.TenantID = house.TenantID.UUID.String()
	}

	return c.JSON(http.StatusOK, formattedHouse)
}

func (app *application) deleteHousesHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid house id"})
	}

	err = app.store.DeleteHouseById(c.Request().Context(), uuid)

	if err != nil {
		slog.Error("error fetching house by id", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}

func (app *application) createHouseHandler(c echo.Context) error {

	var input struct {
		Location  string `json:"location" validate:"required"`
		Block     string `json:"block" validate:"required,len=1"`
		Partition int16  `json:"partition" validate:"required,min=1,max=9"`
		Price     int32  `json:"amount" validate:"required"`
		Occupied  bool   `json:"occupied"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	if err := app.validator.Struct(input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
	}

	_, err := app.store.CreateHouse(c.Request().Context(), db.CreateHouseParams{
		Location:  input.Location,
		Block:     input.Block,
		Partition: input.Partition,
		Price:     input.Price,
		Occupied:  input.Occupied})

	if err != nil {
		slog.Error("error creating house", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (app *application) updateHouseHandler(c echo.Context) error {

	uuid, err := db.ReadUUIDParam(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid house id"})
	}

	house, err := app.store.GetHouseById(c.Request().Context(), uuid)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			slog.Error("error fetching house by id", "err", err)
			return c.JSON(http.StatusNotFound, envelope{"error": "house not found"})

		default:
			slog.Error("error fetching house by id", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	var input struct {
		Location  *string `json:"location"`
		Block     *string `json:"block"`
		Partition *int16  `json:"partition"`
		Price     *int32  `json:"price"`
		Occupied  *bool   `json:"occupied"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	if input.Location != nil {
		house.Location = *input.Location
	}

	if input.Block != nil {
		house.Block = *input.Block
	}

	if input.Partition != nil {
		house.Partition = *input.Partition
	}

	if input.Occupied != nil {
		house.Occupied = *input.Occupied
	}

	if input.Price != nil {
		house.Price = *input.Price
	}

	args := db.UpdateHouseByIdParams{
		ID:        house.HouseID,
		Occupied:  house.Occupied,
		Price:     house.Price,
		Location:  house.Location,
		Block:     house.Block,
		Partition: house.Partition,
		Version:   house.Version,
	}

	err = app.store.UpdateHouseById(c.Request().Context(), args)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.JSON(http.StatusConflict, envelope{"error": "unable to complete request due to an edit conflict"})
		default:
			slog.Error("error updating house ", "err", err)
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}
	}

	return c.JSON(http.StatusOK, nil)
}

func (app *application) bulkHousesHandler(c echo.Context) error {

	var input []struct {
		Location  string   `json:"location" validate:"required"`
		Block     []string `json:"block" validate:"required,min=1,max=5"`
		Partition [][]int  `json:"partition" validate:"gt=0,dive,min=1,max=9,dive,min=1,max=9"`
	}

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, envelope{"error": "invalid request payload"})
	}

	for _, house := range input {
		if err := app.validator.Struct(house); err != nil {
			return c.JSON(http.StatusBadRequest, envelope{"error": err.Error()})
		}
	}

	var housesBulk []db.HouseBulk

	for _, house := range input {
		for i, block := range house.Block {
			for _, pt := range house.Partition[i] {
				housesBulk = append(housesBulk, db.HouseBulk{
					Location:  house.Location,
					Block:     block,
					Partition: pt,
					Occupied:  false,
				})
			}
		}
	}

	err := app.store.BulkInsert(c.Request().Context(), housesBulk)

	if err != nil {
		slog.Error("error bulk inserting houses", "err", err)
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, nil)
}
