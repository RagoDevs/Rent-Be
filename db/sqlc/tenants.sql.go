// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: tenants.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createTenant = `-- name: CreateTenant :exec
INSERT INTO TENANTS (first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

type CreateTenantParams struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
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
		arg.FirstName,
		arg.LastName,
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
SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM
tenants
WHERE tenant_id = $1
`

type GetTenantByIdRow struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HouseID        uuid.UUID `json:"house_id"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
}

func (q *Queries) GetTenantById(ctx context.Context, tenantID uuid.UUID) (GetTenantByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getTenantById, tenantID)
	var i GetTenantByIdRow
	err := row.Scan(
		&i.TenantID,
		&i.FirstName,
		&i.LastName,
		&i.HouseID,
		&i.Phone,
		&i.PersonalIDType,
		&i.PersonalID,
		&i.Active,
		&i.Sos,
		&i.Eos,
	)
	return i, err
}

const getTenants = `-- name: GetTenants :many
SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM tenants
`

type GetTenantsRow struct {
	TenantID       uuid.UUID `json:"tenant_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HouseID        uuid.UUID `json:"house_id"`
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
			&i.TenantID,
			&i.FirstName,
			&i.LastName,
			&i.HouseID,
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
UPDATE tenants 
SET first_name = $1, last_name = $2 ,house_id = $3, phone = $4 ,personal_id_type = $5 ,personal_id = $6 ,active = $7, sos=$8 ,eos = $9
WHERE tenant_id = $10
`

type UpdateTenantParams struct {
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	HouseID        uuid.UUID `json:"house_id"`
	Phone          string    `json:"phone"`
	PersonalIDType string    `json:"personal_id_type"`
	PersonalID     string    `json:"personal_id"`
	Active         bool      `json:"active"`
	Sos            time.Time `json:"sos"`
	Eos            time.Time `json:"eos"`
	TenantID       uuid.UUID `json:"tenant_id"`
}

func (q *Queries) UpdateTenant(ctx context.Context, arg UpdateTenantParams) error {
	_, err := q.db.ExecContext(ctx, updateTenant,
		arg.FirstName,
		arg.LastName,
		arg.HouseID,
		arg.Phone,
		arg.PersonalIDType,
		arg.PersonalID,
		arg.Active,
		arg.Sos,
		arg.Eos,
		arg.TenantID,
	)
	return err
}
