-- name: CreateToken :exec
INSERT INTO token (hash, id, expiry, scope) VALUES ($1, $2, $3, $4);

-- name: DeleteAllToken :exec
DELETE FROM token WHERE scope = $1 AND id = $2;
