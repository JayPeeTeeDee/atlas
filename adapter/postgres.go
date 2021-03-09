package adapter

import (
	"database/sql"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresAdapter struct {
	// TODO: Add connection details
	conn *sql.DB
}

func (p *PostgresAdapter) Connect(dsn string) error {
	// TODO: construct dsn from struct
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	} else {
		p.conn = conn
		return nil
	}
}

func (p *PostgresAdapter) Disconnect() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

func (p *PostgresAdapter) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = p.conn.Query(query, args...)
	return
}

func (p *PostgresAdapter) Exec(query string, args ...interface{}) (result sql.Result, err error) {
	result, err = p.conn.Exec(query, args...)
	return
}

func (p PostgresAdapter) Placeholder() PlaceholderStyle {
	return DollarPlaceholder
}

func (p PostgresAdapter) SpatialType() SpatialExtension {
	return PostGisExtension
}

func (p PostgresAdapter) DatabaseType() DbType {
	return PostgreSQL
}
