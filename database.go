package atlas

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
)

type DatabaseType string

var ErrNoSchema = errors.New("no such schema registered")

const (
	DBType_Postgres DatabaseType = "Postgres"
)

func (dt DatabaseType) IsValid() error {
	switch dt {
	case DBType_Postgres:
		return nil
	}
	return errors.New("Invalid database type")
}

type Database struct {
	databaseType DatabaseType
	adapter      adapter.Adapter
	schemas      map[string]model.Schema
}

func ConnectWithDSN(dbType DatabaseType, dsn string) (*Database, error) {
	var db_adapter adapter.Adapter
	switch dbType {
	case DBType_Postgres:
		db_adapter = &adapter.PostgresAdapter{}
		err := db_adapter.Connect(dsn)
		if err != nil {
			return nil, err
		}
	}
	return &Database{databaseType: dbType, adapter: db_adapter, schemas: make(map[string]model.Schema)}, nil
}

func (d *Database) Disconnect() error {
	return d.adapter.Disconnect()
}

func (d *Database) CreateTable(schemaName string, ifNotExists bool) error {
	schema, ok := d.schemas[schemaName]
	if !ok {
		return errors.New("No such schema registered: " + schemaName)
	}
	sql := query.CompileTableCreation(NewQuery(schema, d), ifNotExists)
	_, err := d.Execute(sql)
	if err == nil {
		statements := query.CompileIndexCreation(NewQuery(schema, d), ifNotExists)
		for _, statement := range statements {
			_, err := d.Execute(statement)
			if err != nil {
				return errors.New("Failed to create indexes for: " + schemaName)
			}
		}
	}
	return err
}

func (d *Database) RegisterModel(target interface{}) error {
	schema, err := model.Parse(target)
	if err != nil {
		return err
	}
	err = schema.SetDefaultValues(target)
	if err != nil {
		return err
	}
	d.schemas[schema.Name] = *schema
	return nil
}

func (d *Database) getSchema(target interface{}) (model.Schema, error) {
	name, err := model.ParseType(target)
	if err != nil {
		return model.Schema{}, err
	}

	return d.GetSchemaByName(name)
}

func (d *Database) GetSchemaByName(name string) (model.Schema, error) {
	if schema, ok := d.schemas[name]; ok {
		return schema, nil
	} else {
		return model.Schema{}, fmt.Errorf("%w: %+v", ErrNoSchema, name)
	}
}

func (d *Database) Model(name string) *Query {
	// Might not be present
	schema := d.schemas[name]
	return NewQuery(schema, d)
}

func (d *Database) Create(object interface{}) (sql.Result, error) {
	schema, err := d.getSchema(object)
	if err != nil {
		return nil, err
	}
	query := NewQuery(schema, d)
	return query.Create(object)
}

func (d *Database) Update(object interface{}) (sql.Result, error) {
	schema, err := d.getSchema(object)
	if err != nil {
		return nil, err
	}
	query := NewQuery(schema, d)
	return query.Update(object)
}

func (d *Database) Execute(query string, args ...interface{}) (sql.Result, error) {
	return d.adapter.Exec(query, args...)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.adapter.Query(query, args...)
}
