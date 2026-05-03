-- name: CreateAPIKey :one
INSERT INTO api_keys (
    user_id,
    name,
    key_hash,
    key_prefix
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListAPIKeysByUser :many
SELECT *
FROM api_keys
WHERE user_id = $1
ORDER BY id DESC;

-- name: GetAPIKeyByHash :one
SELECT *
FROM api_keys
WHERE key_hash = $1
  AND status = 'active';

-- name: UpdateAPIKeyLastUsed :exec
UPDATE api_keys
SET last_used_at = NOW()
WHERE id = $1;


-- name: UpdateAPIKeyStatus :one
UPDATE api_keys
SET status = $2
WHERE id = $1
  AND user_id = $3
RETURNING *;
