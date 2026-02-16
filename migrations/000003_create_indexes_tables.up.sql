CREATE INDEX IF NOT EXISTS idx_operations_wallet_id
    ON operations(wallet_id);

CREATE INDEX IF NOT EXISTS idx_operations_created_at
    ON operations(created_at);