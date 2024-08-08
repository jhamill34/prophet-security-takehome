CREATE TABLE IF NOT EXISTS nodes (
    id  SERIAL PRIMARY KEY,
    ip_addr INET NOT NULL,
    source_id INT NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    version BIGINT DEFAULT 0
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_nodes_unique_ip_source ON nodes (ip_addr, source_id);
