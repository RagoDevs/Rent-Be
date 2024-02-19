-- name: CreateToken :exec
INSERT INTO tokens (hash, admin_id, expiry, scope) VALUES ($1, $2, $3, $4);

-- name: DeleteAllToken :exec
DELETE FROM tokens WHERE scope = $1 AND admin_id = $2;
