-- name: ListAllExistingNodes :many
SELECT DISTINCT n.ip_addr
FROM nodes n
INNER JOIN sources s ON s.id = n.source_id
WHERE s.version < n.version and (n.ip_addr > $1)
ORDER BY n.ip_addr
LIMIT $2;

-- name: BatchInsertNodes :batchexec
INSERT INTO nodes (ip_addr, source_id, version) 
VALUES ($1, $2, $3) 
ON CONFLICT(ip_addr, source_id) 
DO UPDATE 
SET version = EXCLUDED.version;

