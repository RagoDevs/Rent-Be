-- name: CreateTenant :exec
INSERT INTO TENANT
(name, house_id, phone, personal_id_type,personal_id, active, sos, eos) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ;

-- name: GetTenantById :one
SELECT * FROM tenant
WHERE id = $1;

-- name: GetTenantByIdWithHouse :one
SELECT t.id, t.name, t.house_id,h.location, h.block, h.partition, 
t.phone, t.personal_id_type,t.personal_id, t.active, t.sos, t.eos, t.version 
FROM tenant t
JOIN house h ON t.house_id = h.id
WHERE t.id = $1;

-- name: GetTenants :many
SELECT id, name, house_id, 
phone, personal_id_type,personal_id, active, sos, eos
FROM tenant;

-- name: UpdateTenant :exec
UPDATE tenant 
SET name = $1, house_id = $2, phone = $3 ,personal_id_type = $4 ,personal_id = $5 ,active = $6, sos=$7 ,eos = $8, version = uuid_generate_v4()
WHERE id = $9 AND version = $10;

-- name: GetHouseByIdWithTenant :one
SELECT h.id,h.location, h.block, h.partition , h.Occupied, t.name, t.id AS tenant_id
FROM tenant t
Join house h ON h.id = t.house_id
WHERE h.id = $1;

