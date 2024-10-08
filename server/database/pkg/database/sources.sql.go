// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: sources.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createSource = `-- name: CreateSource :one
INSERT INTO sources (name, url, period) 
VALUES ($1, $2, $3) 
ON CONFLICT(name) 
DO NOTHING
RETURNING id, name, url, period, last_execution, version, running
`

type CreateSourceParams struct {
	Name   string
	Url    string
	Period pgtype.Interval
}

func (q *Queries) CreateSource(ctx context.Context, arg CreateSourceParams) (Source, error) {
	row := q.db.QueryRow(ctx, createSource, arg.Name, arg.Url, arg.Period)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Period,
		&i.LastExecution,
		&i.Version,
		&i.Running,
	)
	return i, err
}

const getSource = `-- name: GetSource :one
SELECT id, name, url, period, last_execution, version, running 
FROM sources
WHERE 1=1
AND id = $1
LIMIT 1
`

func (q *Queries) GetSource(ctx context.Context, id int32) (Source, error) {
	row := q.db.QueryRow(ctx, getSource, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Period,
		&i.LastExecution,
		&i.Version,
		&i.Running,
	)
	return i, err
}

const listAllSources = `-- name: ListAllSources :many
SELECT id, name, url, period, last_execution, version, running 
FROM sources
WHERE 1=1
AND id > $1
ORDER BY id
LIMIT $2
`

type ListAllSourcesParams struct {
	ID    int32
	Limit int32
}

func (q *Queries) ListAllSources(ctx context.Context, arg ListAllSourcesParams) ([]Source, error) {
	rows, err := q.db.Query(ctx, listAllSources, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Source
	for rows.Next() {
		var i Source
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Period,
			&i.LastExecution,
			&i.Version,
			&i.Running,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listEligableSources = `-- name: ListEligableSources :many
SELECT id, name, url, period, last_execution, version, running
FROM sources
WHERE 1=1
AND (last_execution IS NULL OR last_execution + period < now())
AND running = TRUE
`

func (q *Queries) ListEligableSources(ctx context.Context) ([]Source, error) {
	rows, err := q.db.Query(ctx, listEligableSources)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Source
	for rows.Next() {
		var i Source
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Period,
			&i.LastExecution,
			&i.Version,
			&i.Running,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const prepareExecution = `-- name: PrepareExecution :one
UPDATE sources
SET last_execution = now(), version = version + 1
WHERE id = $1
RETURNING id, name, url, period, last_execution, version, running
`

func (q *Queries) PrepareExecution(ctx context.Context, id int32) (Source, error) {
	row := q.db.QueryRow(ctx, prepareExecution, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Period,
		&i.LastExecution,
		&i.Version,
		&i.Running,
	)
	return i, err
}

const startSource = `-- name: StartSource :one
UPDATE sources 
SET running = TRUE
WHERE id = $1
RETURNING id, name, url, period, last_execution, version, running
`

func (q *Queries) StartSource(ctx context.Context, id int32) (Source, error) {
	row := q.db.QueryRow(ctx, startSource, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Period,
		&i.LastExecution,
		&i.Version,
		&i.Running,
	)
	return i, err
}

const stopSource = `-- name: StopSource :one
UPDATE sources 
SET running = FALSE, version = version + 1
WHERE id = $1
RETURNING id, name, url, period, last_execution, version, running
`

func (q *Queries) StopSource(ctx context.Context, id int32) (Source, error) {
	row := q.db.QueryRow(ctx, stopSource, id)
	var i Source
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.Period,
		&i.LastExecution,
		&i.Version,
		&i.Running,
	)
	return i, err
}
