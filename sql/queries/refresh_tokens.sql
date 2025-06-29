-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expired_at, revoked_at, user_id)
VALUES (
    $1,
		NOW(),
		NOW(),
		NOW() + INTERVAL '60 days',
		NULL,
		$2
)
RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens 
WHERE token = $1 LIMIT 1;
