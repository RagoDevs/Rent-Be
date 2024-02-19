-- name: CreateTenant :exec
INSERT INTO TENANTS (first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

-- name: GetTenantById :one
SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM
tenants
WHERE tenant_id = $1;

-- name: GetTenants :many
SELECT tenant_id, first_name, last_name, house_id, phone, personal_id_type,personal_id,active,sos,eos FROM tenants;

-- name: UpdateTenant :exec
UPDATE tenants 
SET first_name = $1, last_name = $2 ,house_id = $3, phone = $4 ,personal_id_type = $5 ,personal_id = $6 ,active = $7, sos=$8 ,eos = $9
WHERE tenant_id = $10;

