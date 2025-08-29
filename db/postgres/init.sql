CREATE TABLE IF NOT EXISTS requests (
    id TEXT PRIMARY KEY,
    "from" VARCHAR(10) NOT NULL,
    "to" VARCHAR(10) NOT NULL,
    amount NUMERIC NOT NULL
    );

CREATE TABLE IF NOT EXISTS responses (
    id TEXT PRIMARY KEY,
    success BOOLEAN NOT NULL,
    terms TEXT,
    privacy TEXT,
    query_id TEXT,
    query_from VARCHAR(10),
    query_to VARCHAR(10),
    query_amount NUMERIC,
    info_timestamp BIGINT,
    info_quote NUMERIC,
    result NUMERIC
);

CREATE TABLE IF NOT EXISTS conversion_logs (
    id TEXT PRIMARY KEY,
    timestamp TIMESTAMPTZ NOT NULL,
    request JSONB,
    response JSONB
);
