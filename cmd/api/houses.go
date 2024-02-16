package main

import (
	"errors"
	"net/http"

	"github.com/Hopertz/rmgmt/internal/data"
	"github.com/labstack/echo/v4"
)

func (app *application) listHousesHandler(c echo.Context) error {
	houses, err := app.models.Houses.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, envelope{"houses": houses})

}

func (app *application) showHousesHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	house, err := app.models.Houses.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, envelope{"error": "house not found"})
		default:
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	return c.JSON(http.StatusOK, envelope{"house": house})
}

func (app *application) createHouseHandler(c echo.Context) error {
	var input struct {
		Location  string `json:"location"`
		Block     string `json:"block"`
		Partition int    `json:"partition"`
		Occupied  bool   `json:"occupied"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	house := &data.House{
		Location:  input.Location,
		Block:     input.Block,
		Partition: input.Partition,
		Occupied:  input.Occupied,
	}

	err := app.models.Houses.Insert(house)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, envelope{"house": house})
}

func (app *application) updateHouseHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	house, err := app.models.Houses.Get(uuid)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return c.JSON(http.StatusNotFound, envelope{"error": "house not found"})

		default:
			return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
		}

	}

	var input struct {
		Occupied bool `json:"occupied"`
	}

	if err := c.Bind(&input); err != nil {
		return err
	}

	house.Occupied = input.Occupied

	err = app.models.Houses.Update(house.HouseId, house.Occupied)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, envelope{"house": house})
}

func (app *application) bulkHousesHandler(c echo.Context) error {

	var houses []data.House
	
	if err := c.Bind(&houses); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, envelope{"error": "invalid request payload"})
	}

	err := app.models.Houses.BulkInsert(houses)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, envelope{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, envelope{"houses": houses})

}

// func (app *application) magicbulkHousesHandler(w http.ResponseWriter, r *http.Request) {

// 	type HouseData struct {
// 		Location  string   `json:"location"`
// 		Block     []string `json:"block"`
// 		Partition [][]int  `json:"partition"`
// 	}

// 	type DBHouse struct {
// 		Location  string `json:"location"`
// 		Block     string `json:"block"`
// 		Partition int    `json:"partition"`
// 		Occupied  bool   `json:"occupied"`
// 	}

// 	house_data := `[

// 	{
// 	  "location": "Chanika",
// 	  "block":["A", "B"],
// 	  "partition" : [[1,2], [1]]
// 	},

// 	{
// 	  "location": "Taliani",
// 	  "block":["A", "B"],
// 	  "partition" : [[1,2], [1]]
// 	 },

// 	{
// 	  "location": "Kivule",
// 	  "block":["A", "B"],
// 	  "partition" : [[1,2,3,4,5,6], [1,2,3,4]]
// 	},

// 	{
// 	  "location": "Machimbo",
// 	  "block":["A", "B", "C", "D"],
// 	  "partition" : [[1,2,3,4,5], [1,2,3,4,5],[1,2,3,4,5,6,7,8], [1,2,3]]
// 	},

// 	{
// 	  "location": "UKonga",
// 	  "block":["A", "B","C","D"],
// 	  "partition" : [[1,2], [1], [1] ,[1,2]]
// 	}

// ]`

// 	var houseDB []HouseData
// 	var dbHouses []DBHouse

// 	err := json.Unmarshal([]byte(house_data), &houseDB)

// 	if err != nil {
// 		fmt.Println("Fuck the error", err)
// 	}

// 	for _, house := range houseDB {

// 		for i, block := range house.Block {

// 			for _, pt := range house.Partition[i] {

// 				dbHouses = append(dbHouses, DBHouse{house.Location, block, pt, false})
// 			}
// 		}
// 	}

// 	err = app.writeJSON(w,http.StatusOK,envelope{"houses": dbHouses}, nil)

// 	if err != nil {
// 		app.serverErrorResponse(w,r,err)
// 	}

// }
