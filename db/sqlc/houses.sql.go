// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: houses.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createHouse = `-- name: CreateHouse :one
INSERT INTO house (location, block, partition, price, occupied) VALUES ($1,$2,$3,$4,$5) RETURNING id
`

type CreateHouseParams struct {
	Location  string `json:"location"`
	Block     string `json:"block"`
	Partition int16  `json:"partition"`
	Price     int32  `json:"price"`
	Occupied  bool   `json:"occupied"`
}

func (q *Queries) CreateHouse(ctx context.Context, arg CreateHouseParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createHouse,
		arg.Location,
		arg.Block,
		arg.Partition,
		arg.Price,
		arg.Occupied,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteHouseById = `-- name: DeleteHouseById :exec
DELETE FROM house WHERE id = $1
`

func (q *Queries) DeleteHouseById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteHouseById, id)
	return err
}

const getHouseById = `-- name: GetHouseById :one
SELECT 
  h.id AS house_id,
  h.location, 
  h.block, 
  h.partition, 
  h.price,
  h.Occupied, 
  t.name, 
  t.id AS tenant_id,
  h.version
FROM house h
LEFT JOIN tenant t ON h.id = t.house_id
WHERE h.id = $1
`

type GetHouseByIdRow struct {
	HouseID   uuid.UUID      `json:"house_id"`
	Location  string         `json:"location"`
	Block     string         `json:"block"`
	Partition int16          `json:"partition"`
	Price     int32          `json:"price"`
	Occupied  bool           `json:"occupied"`
	Name      sql.NullString `json:"name"`
	TenantID  uuid.NullUUID  `json:"tenant_id"`
	Version   uuid.UUID      `json:"version"`
}

func (q *Queries) GetHouseById(ctx context.Context, id uuid.UUID) (GetHouseByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getHouseById, id)
	var i GetHouseByIdRow
	err := row.Scan(
		&i.HouseID,
		&i.Location,
		&i.Block,
		&i.Partition,
		&i.Price,
		&i.Occupied,
		&i.Name,
		&i.TenantID,
		&i.Version,
	)
	return i, err
}

const getHouses = `-- name: GetHouses :many
SELECT id,location, block, partition, price , occupied FROM house
`

type GetHousesRow struct {
	ID        uuid.UUID `json:"id"`
	Location  string    `json:"location"`
	Block     string    `json:"block"`
	Partition int16     `json:"partition"`
	Price     int32     `json:"price"`
	Occupied  bool      `json:"occupied"`
}

func (q *Queries) GetHouses(ctx context.Context) ([]GetHousesRow, error) {
	rows, err := q.db.QueryContext(ctx, getHouses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetHousesRow{}
	for rows.Next() {
		var i GetHousesRow
		if err := rows.Scan(
			&i.ID,
			&i.Location,
			&i.Block,
			&i.Partition,
			&i.Price,
			&i.Occupied,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateHouseById = `-- name: UpdateHouseById :exec
UPDATE house
SET location = $1, block = $2, partition = $3, occupied = $4, price = $5, 
version = uuid_generate_v4()
WHERE id = $6 AND version = $7
`

type UpdateHouseByIdParams struct {
	Location  string    `json:"location"`
	Block     string    `json:"block"`
	Partition int16     `json:"partition"`
	Occupied  bool      `json:"occupied"`
	Price     int32     `json:"price"`
	ID        uuid.UUID `json:"id"`
	Version   uuid.UUID `json:"version"`
}

func (q *Queries) UpdateHouseById(ctx context.Context, arg UpdateHouseByIdParams) error {
	_, err := q.db.ExecContext(ctx, updateHouseById,
		arg.Location,
		arg.Block,
		arg.Partition,
		arg.Occupied,
		arg.Price,
		arg.ID,
		arg.Version,
	)
	return err
}
