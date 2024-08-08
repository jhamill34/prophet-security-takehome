-- name: ListAllLists :many 
SELECT *
FROM allowlist
WHERE 1=1
AND id > $1
ORDER BY id
LIMIT $2;

-- name: ListEntriesForAllowList :many 
SELECT *
FROM allowlist_entry 
WHERE 1=1
AND list_id = $1
ORDER BY cidr;

-- name: CreateAllowList :one
INSERT INTO allowlist (name)
VALUES ($1)
RETURNING *;

-- name: DeleteAllowList :exec
DELETE FROM allowlist
WHERE id = $1;

-- name: AddToAllowlist :one
INSERT INTO allowlist_entry (cidr, list_id) 
VALUES ($1, $2)
ON CONFLICT (cidr, list_id) 
DO NOTHING
RETURNING *;

-- name: RemoveFromAllowlist :exec
DELETE FROM allowlist_entry 
WHERE 1=1
AND id = $1 
AND list_id = $2;
