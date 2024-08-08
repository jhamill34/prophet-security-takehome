-- name: ListEntriesForAllowList :many 
SELECT DISTINCT ip_addr
FROM allowlist_entry 
WHERE list_id = $1
ORDER BY ip_addr;

-- name: CreateAllowList :one
INSERT INTO allowlist (name)
VALUES ($1)
RETURNING *;

-- name: DeleteAllowList :exec
DELETE FROM allowlist
WHERE id = $1;

-- name: AddToAllowlist :one
INSERT INTO allowlist_entry (ip_addr, list_id) 
VALUES ($1, $2)
ON CONFLICT (ip_addr, list_id) 
DO NOTHING
RETURNING *;

-- name: RemoveFromAllowlist :exec
DELETE FROM allowlist_entry 
WHERE ip_addr = $1
AND list_id = $2;
