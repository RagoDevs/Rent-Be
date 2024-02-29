-- name: CreatePayment :exec
INSERT INTO payment (tenant_id, amount, start_date, end_date, created_by) VALUES ($1,$2,$3,$4,$5) RETURNING id;


-- name: GetPaymentById :one
SELECT * FROM payment
WHERE id = $1;

-- name: GetDetailedPaymentById :one
SELECT p.id, CONCAT(t.first_name || ' ' || t.last_name) AS tenant_name,
t.id AS tenant_id,p.amount, p.start_date, p.end_date, a.phone AS admin_phone, h.location, h.block, h.partition, 
p.created_at , p.updated_at, p.version  
FROM payment p
JOIN tenant t ON p.tenant_id = t.id
JOIN house h ON t.house_id = h.id
JOIN admin a ON p.created_by = a.id
WHERE p.id = $1;


-- name: GetAllPayments :many
SELECT * FROM payment;


-- name: UpdatePayment :exec
UPDATE payment
SET amount = $1, start_date = $2, end_date = $3, version = uuid_generate_v4(), updated_at = NOW()
WHERE id = $4 AND version = $5;


-- name: DeletePayment :exec    
DELETE FROM payment WHERE id = $1;
