-- name: ListAllNodes :many
SELECT n.ip_addr, n.source_id, n.version, s.last_execution
FROM nodes n
INNER JOIN sources s ON s.id = n.source_id
WHERE 1=1 
AND s.version < n.version 
AND n.ip_addr > $1
ORDER BY n.ip_addr
LIMIT $2;

-- name: ListSourcesNodes :many
SELECT n.ip_addr, n.source_id, n.version, s.last_execution
FROM nodes n
INNER JOIN sources s ON s.id = n.source_id
WHERE 1=1 
AND s.version < n.version 
AND s.id = $3
AND n.ip_addr > $1
ORDER BY n.ip_addr
LIMIT $2;

-- name: ListNodesWithoutAllowlist :many
SELECT n.ip_addr, n.source_id, n.version, s.last_execution
FROM nodes n
INNER JOIN sources s ON s.id = n.source_id
WHERE 1=1
AND s.version < n.version 
AND n.ip_addr > $1
AND NOT n.ip_addr <<= ANY (
    SELECT a.cidr
    FROM allowlist_entry a 
    WHERE 1=1 
    AND a.list_id = $3
)
ORDER BY n.ip_addr
LIMIT $2;

-- name: ListFilteredAllowlistNodes :many
SELECT n.ip_addr, n.source_id, n.version, s.last_execution
FROM nodes n
INNER JOIN sources s ON s.id = n.source_id
WHERE 1=1
AND s.version < n.version 
AND n.ip_addr > $1
AND n.ip_addr <<= ANY (
    SELECT a.cidr
    FROM allowlist_entry a 
    WHERE 1=1 
    AND a.list_id = $3
)
ORDER BY n.ip_addr
LIMIT $2;

-- name: BatchInsertNodes :batchexec
INSERT INTO nodes (ip_addr, source_id, version) 
VALUES ($1, $2, $3) 
ON CONFLICT(ip_addr, source_id) 
DO UPDATE 
SET version = EXCLUDED.version;

