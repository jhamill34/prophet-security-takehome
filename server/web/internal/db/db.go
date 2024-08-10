package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func NewDatabase(ctx context.Context, connection string) *pgx.Conn {
	db, err := pgx.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	return db
}
