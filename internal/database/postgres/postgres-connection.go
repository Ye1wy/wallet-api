package postgres

import (
	"context"
	"log/slog"
	"wallet-api/internal/logger"

	"github.com/jackc/pgx/v5"
)

func Connect(url string, log *slog.Logger) (*pgx.Conn, error) {
	op := "database.postgres.postgres-connection.Connect"

	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Error("Cannot connect to database: ", logger.Err(err), "op", op)
		return nil, err
	}

	return conn, nil
}
