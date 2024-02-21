-- name: CreateTenant :exec
INSERT INTO TENANT
(first_name, last_name, house_id, phone, personal_id_type,personal_id, active, sos, eos) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetTenantById :one
SELECT id, first_name, last_name, house_id, 
phone, personal_id_type,personal_id, active, sos, eos, version 
FROM tenant
WHERE id = $1;

-- name: GetTenants :many
SELECT id, first_name, last_name, house_id, 
phone, personal_id_type,personal_id, active, sos, eos 
FROM tenant;

-- name: UpdateTenant :exec
UPDATE tenant 
SET first_name = $1, last_name = $2 ,house_id = $3, phone = $4 ,personal_id_type = $5 ,personal_id = $6 ,active = $7, sos=$8 ,eos = $9, version = uuid_generate_v4()
WHERE id = $10 AND version = $11;

