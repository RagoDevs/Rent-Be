-- name: CreateHouse :one
INSERT INTO house (location, block, partition, occupied) VALUES ($1,$2,$3,$4) RETURNING id;

-- name: GetHouses :many
SELECT id,location, block, partition , occupied FROM house;

-- name: UpdateHouseById :exec
UPDATE house
SET location = $1, block = $2, partition = $3, occupied = $4, 
version = uuid_generate_v4()
WHERE id = $5 AND version = $6;

-- name: GetHouseById :one
SELECT id,location, block, partition , Occupied, version FROM house WHERE id = $1;






