package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresConnection(DSN string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", DSN)

	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение к БД: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("не удалось пингануть БД: %w", err)
	}

	return conn, nil
}
