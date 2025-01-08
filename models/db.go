package models

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Connect(connString string) *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	return conn
}