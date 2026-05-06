-- name: CreateBalanceTransaction :one
INSERT INTO balance_transactions (
    user_id,
    type,
    amount_cents,
    balance_before_cents,
    balance_after_cents,
    order_id,
    usage_log_id,
    remark
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListBalanceTransactionsByUser :many
SELECT *
FROM balance_transactions
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;
