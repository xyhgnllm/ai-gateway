-- name: CreateUsageLog :one
INSERT INTO usage_logs (
    user_id,
    api_key_id,
    model,
    endpoint,
    status_code,
    success,
    request_body_bytes,
    response_body_bytes,
    cost_cents
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: ListUsageLogsByUser :many
SELECT *
FROM usage_logs
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;
