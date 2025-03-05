-- name: GetAdminByPhone :one
SELECT id, created_at, phone, password_hash, activated, version
FROM admin
WHERE phone = $1;


-- name: CreateAdmin :one
INSERT INTO admin (phone, password_hash, activated)
VALUES ($1, $2, $3 )
RETURNING id, created_at, version;


-- name: UpdateAdmin :one
UPDATE admin
SET phone = $1, password_hash = $2, activated = $3, version = uuid_generate_v4()
WHERE id = $4 AND version = $5
RETURNING version;


-- name: GetHashTokenForAdmin :one
SELECT admin.id, admin.created_at,admin.phone, admin.password_hash,admin.version, admin.activated
FROM admin
INNER JOIN token
ON admin.id = token.id
WHERE token.hash = $1
AND token.scope = $2
AND token.expiry > $3;