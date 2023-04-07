package main

import (
	"encoding/json"
	"fmt"
)

type HouseData struct {
	Location  string     `json:"location"`
	Block     []string   `json:"block"`
	Partition [][]string `json:"partition"`
}

type DBHouse struct {
	Location  string `json:"location"`
	Block     string `json:"block"`
	Partition string `json:"partition"`
}

func main() {
	house_data := `[

	{
	  "location": "Chanika",
	  "block":["A", "B"], 
	  "partition" : [["1","2"], ["1"]]
	},
	
	{
	  "location": "Taliani",
	  "block":["A", "B"],  
	  "partition" : [["1","2"], ["1"]]
	 },
	
	{
	  "location": "Kivule",
	  "block":["A", "B"],  
	  "partition" : [["1","2","3","4","5","6"], ["1","2","3","4"]]
	},
	
	{
	  "location": "Machimbo",
	  "block":["A", "B", "C", "D"],  
	  "partition" : [["1","2","3","4","5"], ["1","2","3","4","5"],["1","2","3","4","5","6","7","8"], ["1","2","3"]]
	},
	
	{
	  "location": "UKonga",
	  "block":["A", "B","C","D"],  
	  "partition" : [["1","2"], ["1"], ["1"] ,["1","2"]]
	}
	
]`

	var houses []HouseData
	var dbHouses []DBHouse

	err := json.Unmarshal([]byte(house_data), &houses)

	if err != nil {
		fmt.Println("Fuck the error", err)
	}

	for _, house := range houses {

		for i, block := range house.Block {

			for _, pt := range house.Partition[i] {

				dbHouses = append(dbHouses, DBHouse{house.Location, block, pt})
			}
		}
	}

	fmt.Println(dbHouses)
	b, err := json.Marshal(dbHouses)

	if err != nil {
		fmt.Println(err)
	}


    fmt.Println(b)
}
