-- name: CreateTenant :exec
INSERT INTO TENANT
(first_name, last_name, house_id, phone, personal_id_type,personal_id, active, sos, eos) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) ;

-- name: GetTenantById :one
SELECT * FROM tenant
WHERE id = $1;

-- name: GetTenantByIdWithHouse :one
SELECT t.id, t.first_name, t.last_name, t.house_id,h.location, h.block, h.partition, 
t.phone, t.personal_id_type,t.personal_id, t.active, t.sos, t.eos, t.version 
FROM tenant t
JOIN house h ON t.house_id = h.id
WHERE t.id = $1;

-- name: GetTenants :many
SELECT id, first_name, last_name, house_id, 
phone, personal_id_type,personal_id, active, sos, eos
FROM tenant;

-- name: UpdateTenant :exec
UPDATE tenant 
SET first_name = $1, last_name = $2 ,house_id = $3, phone = $4 ,personal_id_type = $5 ,personal_id = $6 ,active = $7, sos=$8 ,eos = $9, version = uuid_generate_v4()
WHERE id = $10 AND version = $11;

-- name: GetHouseByIdWithTenant :one
SELECT h.id,h.location, h.block, h.partition , h.Occupied, 
CONCAT(t.first_name || ' ' || t.last_name) AS tenant_name, t.id AS tenant_id
FROM tenant t
Join house h ON h.id = t.house_id
WHERE h.id = $1;

