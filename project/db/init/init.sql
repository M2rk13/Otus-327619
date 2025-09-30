CREATE TABLE IF NOT EXISTS conversion_history (
    id BIGSERIAL PRIMARY KEY,
    currency_from VARCHAR(3) NOT NULL,
    currency_to VARCHAR(3) NOT NULL,
    amount NUMERIC(20, 8) NOT NULL,
    result NUMERIC(20, 8) NOT NULL,
    rate NUMERIC(20, 8) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
