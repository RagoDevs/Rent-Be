-- name: CreateHouse :one
INSERT INTO house (location, block, partition, occupied) VALUES ($1,$2,$3,$4) RETURNING id;

-- name: GetHouses :many
SELECT id,location, block, partition , occupied FROM house;

-- name: UpdateHouseById :exec
UPDATE house
SET location = $1, block = $2, partition = $3, occupied = $4, 
version = uuid_generate_v4(), occupiedBy = $5
WHERE id = $6 AND version = $7;

-- name: GetHouseById :one
SELECT id,location, block, partition , Occupied, occupiedBy, version FROM house WHERE id = $1;

-- name: GetHouseByIdWithTenant :one
SELECT h.id,h.location, h.block, h.partition , h.Occupied, 
CONCAT(t.first_name || ' ' || t.last_name) AS tenant_name, t.id AS tenant_id, h.version 
FROM house h
Join tenant t ON h.occupiedBy = t.id
WHERE h.id = $1;





