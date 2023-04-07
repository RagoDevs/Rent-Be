package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type House struct {
	HouseId   string `json:"house_id"`
	Location  string `json:"location"`
	Block     string `json:"block"`
	Partition int `json:"partition"`
	Occupied  bool   `json:"occupied"`
}

type HouseModel struct {
	DB *sql.DB
}

func (h HouseModel) Insert(house *House) error {
	query := `INSERT INTO houses (location,block,partition, Occupied) VALUES ($1,$2,$3,$4) RETURNING house_id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	args := []interface{}{house.Location, house.Block, house.Partition, house.Occupied}

	err := h.DB.QueryRowContext(ctx, query, args...).Scan(&house.HouseId)

	if err != nil {
		return err
	}

	return nil

}

func (h HouseModel) GetAll() ([]*House, error) {
	query := `SELECT house_id,location, block, partition , occupied FROM houses`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := h.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	houses := []*House{}

	for rows.Next() {
		var house House

		err := rows.Scan(
			&house.HouseId,
			&house.Location,
			&house.Block,
			&house.Partition,
			&house.Occupied,
		)

		if err != nil {
			return nil, err
		}

		houses = append(houses, &house)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return houses, nil

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

func (h HouseModel) Update(house *House) error {
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
