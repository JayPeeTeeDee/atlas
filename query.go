package atlas

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
	"github.com/georgysavva/scany/dbscan"
)

type Query struct {
	// TODO: Change to take in model
	schema   model.Schema
	builder  *query.Builder
	database *Database
	compiler *query.Compiler

	buildErrors []error
}

type Result struct {
	Error        error
	RowsAffected int
}

func NewQuery(schema model.Schema, database *Database) *Query {
	return &Query{
		schema:   schema,
		builder:  query.NewBuilder(),
		database: database,
		compiler: &query.Compiler{
			AdapterInfo: database.adapter,
			Schema:      schema,
		},
		buildErrors: make([]error, 0),
	}
}

func (q *Query) Error() error {
	if len(q.buildErrors) == 0 {
		return nil
	} else {
		errList := make([]string, len(q.buildErrors))
		for i, err := range q.buildErrors {
			errList[i] = err.Error()
		}
		return errors.New(strings.Join(errList, ","))
	}
}

/* Functions for building up query */
func (q *Query) Select(columns ...string) *Query {
	missingCols := q.getMissingCols(columns...)
	if len(missingCols) > 0 {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Missing cols: %s", strings.Join(missingCols, ",")))
	}
	q.builder.Selections.AddAll(columns...)
	return q
}

func (q *Query) Omit(columns ...string) *Query {
	missingCols := q.getMissingCols(columns...)
	if len(missingCols) > 0 {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Missing cols: %s", strings.Join(missingCols, ",")))
	}
	q.builder.Omissions.AddAll(columns...)
	return q
}

func (q *Query) Where(clause query.Clause) *Query {
	if !clause.IsValid(q.schema.FieldsByName) {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Invalid clause of type: %s", clause.Condition()))
	}
	q.builder.Where(clause)
	return q
}

func (q *Query) Limit(count uint64) *Query {
	q.builder.Limit = count
	return q
}

func (q *Query) Offset(count uint64) *Query {
	q.builder.Offset = count
	return q
}

func (q *Query) OrderBy(column string, desc bool) *Query {
	missingCols := q.getMissingCols(column)
	if len(missingCols) > 0 {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Mising cols: %s", strings.Join(missingCols, ",")))
	}
	q.builder.OrderBy(column, desc)
	return q
}

func (q *Query) CoveredBy(target model.SpatialObject) *Query {
	if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.schema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.CoveredBy{Column: q.schema.LocationFieldNames.Keys()[0], Target: target})
	} else {
		q.builder.Where(query.CoveredBy{Column: q.schema.RegionFieldNames.Keys()[0], Target: target})
	}
	return q
}

func (q *Query) Covers(target model.SpatialObject) *Query {
	if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.schema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.Covers{Column: q.schema.LocationFieldNames.Keys()[0], Target: target})
	} else {
		q.builder.Where(query.Covers{Column: q.schema.RegionFieldNames.Keys()[0], Target: target})
	}
	return q
}

func (q *Query) WithinRangeOf(targets []model.SpatialObject, rangeMeters float64) *Query {
	if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.schema.LocationFieldNames.Size() == 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else {
		q.builder.Where(query.WithinRangeOf{Column: q.schema.RegionFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
	}
	return q
}

func (q *Query) HasWithinRange(targets []model.SpatialObject, rangeMeters float64) *Query {
	if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.schema.LocationFieldNames.Size()+q.schema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.schema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.HasWithinRange{Column: q.schema.LocationFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
	} else {
		q.builder.Where(query.HasWithinRange{Column: q.schema.RegionFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
	}
	return q
}

/* Functions for execution of query */

/* SELECT STATEMENTS */
func (q *Query) Count(count *int) error {
	if q.Error() != nil {
		return q.Error()
	}
	q.builder.QueryType = query.SelectQuery
	// TODO: order by primary key, limit 1
	statement, args := q.compiler.CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	if err != nil {
		return err
	}
	return dbscan.ScanOne(count, rows)
}

func (q *Query) First(response interface{}) error {
	if q.Error() != nil {
		return q.Error()
	}
	q.builder.QueryType = query.SelectQuery
	// TODO: order by primary key, limit 1
	q.builder.Limit = 1
	statement, args := q.compiler.CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	if err != nil {
		return err
	}
	return dbscan.ScanOne(response, rows)
}

func (q *Query) All(response interface{}) error {
	if q.Error() != nil {
		return q.Error()
	}
	q.builder.QueryType = query.SelectQuery
	statement, args := q.compiler.CompileSQL(*q.builder)
	rows, err := q.database.Query(statement, args...)
	// TODO: return wrapped error
	if err != nil {
		return err
	}

	return dbscan.ScanAll(response, rows)
}

func (q *Query) Create(object interface{}) (sql.Result, error) {
	if q.Error() != nil {
		return nil, q.Error()
	}
	q.builder.QueryType = query.InsertQuery
	vals, err := model.ParseObject(object, q.schema)
	if err != nil {
		return nil, err
	}
	q.builder.InsertValues = vals
	statement, args := q.compiler.CompileSQL(*q.builder)
	// TODO: Make use of result
	return q.database.Execute(statement, args...)
}

func (q *Query) Update(object interface{}) (sql.Result, error) {
	if q.Error() != nil {
		return nil, q.Error()
	}
	q.builder.QueryType = query.UpdateQuery
	vals, err := model.ParseObject(object, q.schema)
	if err != nil {
		return nil, err
	}
	if len(vals) > 1 {
		return nil, errors.New("Can only update 1 record each time")
	}
	q.builder.InsertValues = vals
	statement, args := q.compiler.CompileSQL(*q.builder)
	// TODO: Make use of result
	return q.database.Execute(statement, args...)
}

func (q *Query) getMissingCols(columns ...string) []string {
	missingCols := make([]string, 0)
	for _, col := range columns {
		if !q.schema.AllFieldNames.Contains(col) {
			missingCols = append(missingCols, col)
		}
	}
	return missingCols
}
