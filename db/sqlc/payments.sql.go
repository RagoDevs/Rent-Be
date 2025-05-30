// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: payments.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPayment = `-- name: CreatePayment :exec
INSERT INTO payment (tenant_id, amount, start_date, end_date, created_by) VALUES ($1,$2,$3,$4,$5) RETURNING id
`

type CreatePaymentParams struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	Amount    int32     `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) error {
	_, err := q.db.ExecContext(ctx, createPayment,
		arg.TenantID,
		arg.Amount,
		arg.StartDate,
		arg.EndDate,
		arg.CreatedBy,
	)
	return err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM payment WHERE id = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePayment, id)
	return err
}

const getAllPayments = `-- name: GetAllPayments :many
SELECT p.id, t.name AS tenant_name,
t.id AS tenant_id,p.amount, p.start_date, p.end_date, a.email AS admin_email, h.location, h.block, h.partition, 
p.created_at , p.updated_at, p.version  
FROM payment p
JOIN tenant t ON p.tenant_id = t.id
JOIN house h ON t.house_id = h.id
JOIN admin a ON p.created_by = a.id
`

type GetAllPaymentsRow struct {
	ID         uuid.UUID `json:"id"`
	TenantName string    `json:"tenant_name"`
	TenantID   uuid.UUID `json:"tenant_id"`
	Amount     int32     `json:"amount"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	AdminEmail string    `json:"admin_email"`
	Location   string    `json:"location"`
	Block      string    `json:"block"`
	Partition  int16     `json:"partition"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    uuid.UUID `json:"version"`
}

func (q *Queries) GetAllPayments(ctx context.Context) ([]GetAllPaymentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllPaymentsRow{}
	for rows.Next() {
		var i GetAllPaymentsRow
		if err := rows.Scan(
			&i.ID,
			&i.TenantName,
			&i.TenantID,
			&i.Amount,
			&i.StartDate,
			&i.EndDate,
			&i.AdminEmail,
			&i.Location,
			&i.Block,
			&i.Partition,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Version,
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

const getDetailedPaymentById = `-- name: GetDetailedPaymentById :one
SELECT p.id, t.name AS tenant_name,
t.id AS tenant_id,p.amount, p.start_date, p.end_date, a.email AS admin_email, h.location, h.block, h.partition, 
p.created_at , p.updated_at, p.version  
FROM payment p
JOIN tenant t ON p.tenant_id = t.id
JOIN house h ON t.house_id = h.id
JOIN admin a ON p.created_by = a.id
WHERE p.id = $1
`

type GetDetailedPaymentByIdRow struct {
	ID         uuid.UUID `json:"id"`
	TenantName string    `json:"tenant_name"`
	TenantID   uuid.UUID `json:"tenant_id"`
	Amount     int32     `json:"amount"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	AdminEmail string    `json:"admin_email"`
	Location   string    `json:"location"`
	Block      string    `json:"block"`
	Partition  int16     `json:"partition"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Version    uuid.UUID `json:"version"`
}

func (q *Queries) GetDetailedPaymentById(ctx context.Context, id uuid.UUID) (GetDetailedPaymentByIdRow, error) {
	row := q.db.QueryRowContext(ctx, getDetailedPaymentById, id)
	var i GetDetailedPaymentByIdRow
	err := row.Scan(
		&i.ID,
		&i.TenantName,
		&i.TenantID,
		&i.Amount,
		&i.StartDate,
		&i.EndDate,
		&i.AdminEmail,
		&i.Location,
		&i.Block,
		&i.Partition,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Version,
	)
	return i, err
}

const getPaymentById = `-- name: GetPaymentById :one
SELECT id, tenant_id, amount, start_date, end_date, version, created_at, created_by, updated_at FROM payment
WHERE id = $1
`

func (q *Queries) GetPaymentById(ctx context.Context, id uuid.UUID) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentById, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.TenantID,
		&i.Amount,
		&i.StartDate,
		&i.EndDate,
		&i.Version,
		&i.CreatedAt,
		&i.CreatedBy,
		&i.UpdatedAt,
	)
	return i, err
}

const updatePayment = `-- name: UpdatePayment :exec
UPDATE payment
SET amount = $1, start_date = $2, end_date = $3, version = uuid_generate_v4(), updated_at = NOW()
WHERE id = $4 AND version = $5
`

type UpdatePaymentParams struct {
	Amount    int32     `json:"amount"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	ID        uuid.UUID `json:"id"`
	Version   uuid.UUID `json:"version"`
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) error {
	_, err := q.db.ExecContext(ctx, updatePayment,
		arg.Amount,
		arg.StartDate,
		arg.EndDate,
		arg.ID,
		arg.Version,
	)
	return err
}
