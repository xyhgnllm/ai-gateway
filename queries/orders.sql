-- name: CreateOrder :one
INSERT INTO orders (
    user_id,
    order_no,
    amount_cents
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ListOrdersByUser :many
SELECT *
FROM orders
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;

-- name: ListOrders :many
SELECT *
FROM orders
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: MarkOrderPaid :one
UPDATE orders
SET status = 'paid',
    paid_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND status = 'pending'
RETURNING *;

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE id = $1;
