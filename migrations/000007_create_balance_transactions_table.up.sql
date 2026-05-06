CREATE TABLE balance_transactions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    amount_cents BIGINT NOT NULL,
    balance_before_cents BIGINT NOT NULL,
    balance_after_cents BIGINT NOT NULL,
    order_id BIGINT REFERENCES orders(id) ON DELETE SET NULL,
    usage_log_id BIGINT REFERENCES usage_logs(id) ON DELETE SET NULL,
    remark TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_balance_transactions_user_id ON balance_transactions(user_id);
CREATE INDEX idx_balance_transactions_order_id ON balance_transactions(order_id);
CREATE INDEX idx_balance_transactions_usage_log_id ON balance_transactions(usage_log_id);
CREATE INDEX idx_balance_transactions_created_at ON balance_transactions(created_at);
