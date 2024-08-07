-- name: CreateSource :one
INSERT INTO sources (name, url, period) 
VALUES ($1, $2, $3) 
ON CONFLICT(name) 
DO NOTHING
RETURNING *;

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
AND name = $1
LIMIT 1;

-- name: PrepareExecution :one
UPDATE sources
SET last_execution = now(), version = version + 1
WHERE name = $1
RETURNING *;

-- name: StopSource :one
UPDATE sources 
SET running = FALSE, version = version + 1
WHERE name = $1
RETURNING *;

-- name: StartSource :one
UPDATE sources 
SET running = TRUE
WHERE name = $1
RETURNING *;
