-- name: CreateSource :one
INSERT INTO sources (name, url, period) 
VALUES ($1, $2, $3) 
ON CONFLICT(name) 
DO NOTHING
RETURNING *;

-- name: ListAllSources :many
SELECT * 
FROM sources
WHERE 1=1
AND id > $1
ORDER BY id
LIMIT $2;

-- name: ListEligableSources :many
SELECT *
FROM sources
WHERE 1=1
AND (last_execution IS NULL OR last_execution + period < now())
AND running = TRUE;  

-- name: GetSource :one
SELECT * 
FROM sources
WHERE 1=1
AND id = $1
LIMIT 1;

-- name: PrepareExecution :one
UPDATE sources
SET last_execution = now(), version = version + 1
WHERE id = $1
RETURNING *;

-- name: StopSource :one
UPDATE sources 
SET running = FALSE, version = version + 1
WHERE id = $1
RETURNING *;

-- name: StartSource :one
UPDATE sources 
SET running = TRUE
WHERE id = $1
RETURNING *;
