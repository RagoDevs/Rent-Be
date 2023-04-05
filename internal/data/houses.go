package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type House struct {
	HouseId   string `json:"house_id"`
	Location  string `json:"location"`
	Block     string `json:"block"`
	Partition string `json:"partition"`
	Occupied  bool   `json:occupied"`
}

type HouseModel struct {
	DB *sql.DB
}

func (h HouseModel) Insert(house *House) error {
	query := fmt.Sprintf(`INSERT INTO houses (house_id,location,block,partition, Occupied) VALUES (%s,$1,$2,$3,$4)`, "uuid_generate_v4()")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{house.Location, house.Block, house.Partition, house.Occupied}

	_, err := h.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil

}

func (h HouseModel) Get(house_id string) (*House, error) {
	query := `SELECT house_id,location, block, partition , Occupied FROM houses
	WHERE house_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var house House

	err := h.DB.QueryRowContext(ctx, query, house_id).Scan(
		&house.HouseId,
		&house.Location,
		&house.Block,
		&house.Partition,
		&house.Occupied,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err

		}
	}

	return &house, nil

}

func (h HouseModel) Update(house House) error {
	query := `UPDATE houses
	SET occupied = $1
	WHERE house_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{house.Occupied, house.HouseId}

	_, err := h.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}
