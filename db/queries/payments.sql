-- name: CreatePayment :exec
INSERT INTO payment (tenant_id, period,start_date, end_date, renewed) VALUES ($1,$2,$3,$4,$5);


-- name: GetPaymentById :one
SELECT id, tenant_id, period, start_date, end_date, renewed , version  FROM payment
WHERE id = $1;


-- name: GetUnrenewedByTenantId :one
SELECT id, tenant_id, period, start_date, end_date, renewed, version FROM payment 
WHERE renewed = false and tenant_id = $1;


-- name: GetAllPayments :many
SELECT id,tenant_id, period, start_date, end_date, renewed, version FROM payment;


-- name: UpdatePayment :exec
UPDATE payment
SET period = $1, start_date = $2, end_date = $3, renewed = $4, version = uuid_generate_v4()
WHERE id = $5 AND version = $6;


-- name: DeletePayment :exec    
DELETE FROM payment WHERE id = $1;
