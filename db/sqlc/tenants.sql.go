// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: tenants.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTenant = `-- name: CreateTenant :exec
INSERT INTO TENANT
(name, house_id, phone, personal_id_type,personal_id, active, sos, eos) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

type CreateTenantParams struct {
	Name           string    `json:"name"`
	HouseID        uuid.UUID `json:"house_id"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
}

func (q *Queries) CreateTenant(ctx context.Context, arg CreateTenantParams) error {
	_, err := q.db.ExecContext(ctx, createTenant,
		arg.Name,
		arg.HouseID,
		arg.Phone,
		arg.PersonalIDType,
		arg.PersonalID,
		arg.Active,
		arg.Sos,
		arg.Eos,
	)
	return err
}

const getTenantById = `-- name: GetTenantById :one
SELECT t.id AS tenant_id, t.name, t.house_id,h.location, h.block, h.partition, h.price ,
t.phone, t.personal_id_type,t.personal_id, t.active, t.sos, t.eos, t.version 
FROM tenant t
JOIN house h ON t.house_id = h.id
WHERE t.id = $1
`

type GetTenantByIdRow struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	Name           string    `json:"name"`
	HouseID        uuid.UUID `json:"house_id"`
	Location       string    `json:"location"`
	Block          string    `json:"block"`
	Partition      int16     `json:"partition"`
	Price          int32     `json:"price"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
	Version        uuid.UUID `json:"version"`
}

func (q *Queries) GetTenantById(ctx context.Context, id uuid.UUID) (GetTenantByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getTenantById, id)
	var i GetTenantByIdRow
	err := row.Scan(
		&i.TenantID,
		&i.Name,
		&i.HouseID,
		&i.Location,
		&i.Block,
		&i.Partition,
		&i.Price,
		&i.Phone,
		&i.PersonalIDType,
		&i.PersonalID,
		&i.Active,
		&i.Sos,
		&i.Eos,
		&i.Version,
	)
	return i, err
}

const getTenants = `-- name: GetTenants :many
SELECT 
    t.id, 
    t.name, 
    h.location ,
    h.block ,
    h.partition,
    h.price,
    t.phone, 
    t.personal_id_type,
    t.personal_id, 
    t.active, 
    t.sos, 
    t.eos
FROM tenant t
JOIN house h ON t.house_id = h.id
`

type GetTenantsRow struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Location       string    `json:"location"`
	Block          string    `json:"block"`
	Partition      int16     `json:"partition"`
	Price          int32     `json:"price"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
}

func (q *Queries) GetTenants(ctx context.Context) ([]GetTenantsRow, error) {
	rows, err := q.db.QueryContext(ctx, getTenants)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetTenantsRow{}
	for rows.Next() {
		var i GetTenantsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Location,
			&i.Block,
			&i.Partition,
			&i.Price,
			&i.Phone,
			&i.PersonalIDType,
			&i.PersonalID,
			&i.Active,
			&i.Sos,
			&i.Eos,
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

const updateTenant = `-- name: UpdateTenant :exec
UPDATE tenant 
SET name = $1, house_id = $2, phone = $3 ,personal_id_type = $4 ,personal_id = $5 ,active = $6, sos=$7 ,eos = $8, version = uuid_generate_v4()
WHERE id = $9 AND version = $10
`

type UpdateTenantParams struct {
	Name           string    `json:"name"`
	HouseID        uuid.UUID `json:"house_id"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
	ID             uuid.UUID `json:"id"`
	Version        uuid.UUID `json:"version"`
}

func (q *Queries) UpdateTenant(ctx context.Context, arg UpdateTenantParams) error {
	_, err := q.db.ExecContext(ctx, updateTenant,
		arg.Name,
		arg.HouseID,
		arg.Phone,
		arg.PersonalIDType,
		arg.PersonalID,
		arg.Active,
		arg.Sos,
		arg.Eos,
		arg.ID,
		arg.Version,
	)
	return err
}
