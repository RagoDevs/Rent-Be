-- name: CreateHouse :one
INSERT INTO house (location, block, partition, price, occupied) VALUES ($1,$2,$3,$4,$5) RETURNING id;

-- name: GetHouses :many
SELECT id,location, block, partition, price , occupied FROM house;

-- name: UpdateHouseById :exec
UPDATE house
SET location = $1, block = $2, partition = $3, occupied = $4, price = $5, 
version = uuid_generate_v4()
WHERE id = $6 AND version = $7;

-- name: GetHouseById :one
SELECT 
  h.id AS house_id,
  h.location, 
  h.block, 
  h.partition, 
  h.price,
  h.Occupied, 
  t.name, 
  t.id AS tenant_id,
  h.version
FROM house h
LEFT JOIN tenant t ON h.id = t.house_id
WHERE h.id = $1;

-- name: DeleteHouseById :exec    
DELETE FROM house WHERE id = $1;





