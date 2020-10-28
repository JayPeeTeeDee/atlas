package adapter

import "database/sql"

type Adapter interface {
	Connect(dsn string) error
	Disconnect() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}
