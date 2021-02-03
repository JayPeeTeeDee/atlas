package atlas

import (
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
	"github.com/georgysavva/scany/dbscan"
)

type Query struct {
	// TODO: Change to take in model
	schema   model.Schema
	builder  *query.Builder
	database *Database
}

type Result struct {
	Error        error
	RowsAffected int
}

func NewQuery(schema model.Schema, database *Database) *Query {
	return &Query{
		schema: schema,
		builder: &query.Builder{
			TableName: schema.Table,
		},
		database: database,
	}
}

/* Functions for building up query */
func (q *Query) Where(clause query.Clause) *Query {
	// TODO: check that table has column
	q.builder.Where(clause)
	return q
}

func (q *Query) Limit(count uint64) *Query {
	// TODO: check more than 0
	q.builder.Limit = count
	return q
}

func (q *Query) Offset(count uint64) *Query {
	// TODO: check more than 0
	q.builder.Offset = count
	return q
}

func (q *Query) OrderBy(column string, desc bool) *Query {
	// TODO: check that table has column
	q.builder.OrderBy(column, desc)
	return q
}

/* Functions for execution of query */

/* SELECT STATEMENTS */
func (q *Query) Count(count *int) error {
	q.builder.QueryType = query.SelectQuery
	// TODO: order by primary key, limit 1
	statement, args := CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	if err != nil {
		return err
	}
	return dbscan.ScanOne(count, rows)
}

func (q *Query) First(response interface{}) error {
	q.builder.QueryType = query.SelectQuery
	// TODO: order by primary key, limit 1
	q.builder.Limit = 1
	statement, args := CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	if err != nil {
		return err
	}
	return dbscan.ScanOne(response, rows)
}

func (q *Query) All(response interface{}) error {
	q.builder.QueryType = query.SelectQuery
	statement, args := CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	// TODO: return wrapped error
	if err != nil {
		return err
	}

	return dbscan.ScanAll(response, rows)
}

/*  */