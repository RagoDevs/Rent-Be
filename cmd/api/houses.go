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
		Partition int    `json:"partition"`
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
		Occupied bool `json:"occupied"`
	}

	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	house.Occupied = input.Occupied

	err = app.models.Houses.Update(house.HouseId, house.Occupied)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"house": house}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) bulkHousesHandler(w http.ResponseWriter, r *http.Request) {

	var houses []data.House
	err := app.readBulKJSON(w, r, &houses)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Houses.BulkInsert(houses)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"houses": houses}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

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
