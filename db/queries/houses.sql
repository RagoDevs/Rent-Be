-- name: CreateHouse :one
INSERT INTO house (location, block, partition, occupied) VALUES ($1,$2,$3,$4) RETURNING id;

-- name: GetHouses :many
SELECT id,location, block, partition , occupied FROM house;

-- name: UpdateHouseById :exec
UPDATE house
SET occupied = $1, version = uuid_generate_v4()
WHERE id = $2 AND version = $3;

-- name: GetHouseById :one
SELECT id,location, block, partition , Occupied FROM house
WHERE id = $1;




