CREATE TABLE IF NOT EXISTS allowlist (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS allowlist_entry (
    id SERIAL PRIMARY KEY,
    ip_addr CIDR NOT NULL,
    list_id INT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_allowlist_entry_ip_addr_list_id ON allowlist_entry(ip_addr, list_id);
