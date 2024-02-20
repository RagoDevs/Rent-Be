-- name: GetAdminByEmail :one
SELECT admin_id, created_at, email, password_hash, activated, version
FROM admins
WHERE email = $1;


-- name: InsertAdmin :one
INSERT INTO admins (email, password_hash, activated)
VALUES ($1, $2, $3 )
RETURNING admin_id, created_at, version;


-- name: UpdateAdmin :one
UPDATE admins
SET email = $1, password_hash = $2, activated = $3, version = version + 1
WHERE admin_id = $4 AND version = $5
RETURNING version;


-- name: GetHashTokenForAdmin :one
SELECT admins.admin_id, admins.created_at,admins.email, admins.password_hash, admins.activated, admins.version
FROM admins
INNER JOIN tokens
ON admins.admin_id = tokens.admin_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3;