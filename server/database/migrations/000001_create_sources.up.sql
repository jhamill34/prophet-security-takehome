CREATE TABLE IF NOT EXISTS sources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    url VARCHAR(2000) NOT NULL,
    period INTERVAL NOT NULL,
    last_execution TIMESTAMP,
    version BIGINT DEFAULT 0,
    running BOOLEAN DEFAULT TRUE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_sources_unique_name ON sources (name);

