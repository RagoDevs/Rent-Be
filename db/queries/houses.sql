-- name: CreateHouse :one
INSERT INTO houses (location,block,partition, occupied) VALUES ($1,$2,$3,$4) RETURNING house_id;

-- name: GetHouses :many
SELECT house_id,location, block, partition , occupied FROM houses;

-- name: UpdateHouseById :exec
UPDATE houses
SET occupied = $1
WHERE house_id = $2;

-- name: GetHouseById :one
SELECT house_id,location, block, partition , Occupied FROM houses
WHERE house_id = $1;


