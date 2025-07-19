-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid (),
		NOW(),
		NOW(),
		$1,
		$2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = $1 LIMIT 1;

-- name: GetUserByRefreshToken :one
SELECT u.*
FROM users u
INNER JOIN refresh_tokens rt ON u.id = rt.user_id
WHERE rt.token = $1;

-- name: UpdateUserByUserId :one
UPDATE users
SET updated_at = NOW(), email = $1, hashed_password = $2  
WHERE id = $3
RETURNING *;

-- name: UpdateUserChirpRedTRUE :one
UPDATE users 
SET  is_chirpy_red = TRUE 
WHERE id = $1
RETURNING *;
