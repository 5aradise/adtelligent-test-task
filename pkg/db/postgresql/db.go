package postgresql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const driverName = "postgres"

func New(dsn string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
