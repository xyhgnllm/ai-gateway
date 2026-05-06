-- name: CountUsers :one
SELECT COUNT(*)::BIGINT
FROM users;

-- name: CountOrders :one
SELECT COUNT(*)::BIGINT
FROM orders;

-- name: SumPaidOrderAmount :one
SELECT COALESCE(SUM(amount_cents), 0)::BIGINT
FROM orders
WHERE status = 'paid';

-- name: CountUsageLogs :one
SELECT COUNT(*)::BIGINT
FROM usage_logs;

-- name: SumUsageCost :one
SELECT COALESCE(SUM(cost_cents), 0)::BIGINT
FROM usage_logs
WHERE success = true;


-- name: ListModelUsageStats :many
SELECT
    model,
    COUNT(*)::BIGINT AS total_count,
    COUNT(*) FILTER (WHERE success = true)::BIGINT AS success_count,
    COALESCE(SUM(cost_cents), 0)::BIGINT AS total_cost_cents
FROM usage_logs
GROUP BY model
ORDER BY total_count DESC;


-- name: ListUserUsageStats :many
SELECT
    user_id,
    COUNT(*)::BIGINT AS total_count,
    COUNT(*) FILTER (WHERE success = true)::BIGINT AS success_count,
    COALESCE(SUM(cost_cents), 0)::BIGINT AS total_cost_cents
FROM usage_logs
GROUP BY user_id
ORDER BY total_cost_cents DESC;


-- name: ListDailyUsageStats :many
SELECT
    DATE(created_at) AS date,
    COUNT(*)::BIGINT AS total_count,
    COUNT(*) FILTER (WHERE success = true)::BIGINT AS success_count,
    COALESCE(SUM(cost_cents), 0)::BIGINT AS total_cost_cents
FROM usage_logs
GROUP BY DATE(created_at)
ORDER BY date DESC;
