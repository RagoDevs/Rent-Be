-- name: CreatePayment :exec
INSERT INTO payments (tenant_id,period,start_date,end_date, renewed) VALUES ($1,$2,$3,$4,$5);


-- name: GetPaymentById :one
SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments 
WHERE payment_id = $1;


-- name: GetUnrenewedByTenantId :one
SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments 
WHERE renewed = false and tenant_id = $1;


-- name: GetAllPayments :many
SELECT payment_id,tenant_id, period, start_date, end_date, renewed FROM payments;


-- name: UpdatePayment :exec
UPDATE payments
SET period = $1, start_date = $2, end_date = $3, renewed = $4
WHERE payment_id = $5;


-- name: DeletePayment :exec    
DELETE FROM payments WHERE payment_id = $1;
