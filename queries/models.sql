-- name: CreateModel :one
INSERT INTO models (
    name,
    provider,
    input_price_per_1k_cents,
    output_price_per_1k_cents
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListModels :many
SELECT *
FROM models
ORDER BY id DESC;

-- name: GetActiveModelByName :one
SELECT *
FROM models
WHERE name = $1
  AND status = 'active';

-- name: UpdateModelStatus :one
UPDATE models
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
