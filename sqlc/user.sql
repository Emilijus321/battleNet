-- queries/user.sql

-- name: GetUserByID :one
SELECT * FROM "user"
WHERE user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM "user"
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM "user"
WHERE username = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY created_at DESC
    LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO "user" (
    email,
    password_hash,
    first_name,
    last_name,
    username,
    role,
    avatar_url
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         )
    RETURNING *;

-- name: UpdateUser :one
UPDATE "user"
SET
    email = $2,
    first_name = $3,
    last_name = $4,
    username = $5,
    avatar_url = $6,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: UpdateUserPassword :one
UPDATE "user"
SET
    password_hash = $2,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE "user"
SET
    last_login_at = NOW()
WHERE user_id = $1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE user_id = $1;

-- name: UpdateUserRole :one
UPDATE "user"
SET
    role = $2,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: UpdateUserActiveStatus :one
UPDATE "user"
SET
    is_active = $2,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: VerifyUserEmail :one
UPDATE "user"
SET
    email_verified = true,
    updated_at = NOW()
WHERE user_id = $1
    RETURNING *;

-- name: SearchUsers :many
SELECT * FROM "user"
WHERE
    username ILIKE $1 OR
    email ILIKE $1 OR
    first_name ILIKE $1 OR
    last_name ILIKE $1
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;

-- name: GetUsersCount :one
SELECT COUNT(*) FROM "user";

-- name: GetUserStats :one
SELECT
    COUNT(*) as total_users,
    COUNT(CASE WHEN email_verified = true THEN 1 END) as verified_users,
    COUNT(CASE WHEN is_active = true THEN 1 END) as active_users,
    COUNT(CASE WHEN role = 'admin' THEN 1 END) as admin_users
FROM "user";