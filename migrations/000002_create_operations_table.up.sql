CREATE TABLE IF NOT EXISTS operations (
    id UUID PRIMARY KEY,
    wallet_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    amount BIGINT NOT NULL CHECK (amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_wallet
      FOREIGN KEY (wallet_id)
          REFERENCES wallets(id)
          ON DELETE CASCADE,

    CONSTRAINT operation_type_check
      CHECK (type IN ('DEPOSIT', 'WITHDRAW'))
);
