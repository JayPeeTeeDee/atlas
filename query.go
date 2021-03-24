package atlas

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
	"github.com/georgysavva/scany/dbscan"
)

type Query struct {
	// TODO: Change to take in model
	mainSchema  model.Schema
	joinSchemas map[string]model.Schema
	builder     *query.Builder
	database    *Database
	buildErrors []error
}

type Result struct {
	Error        error
	RowsAffected int
}

func NewQuery(schema model.Schema, database *Database) *Query {
	return &Query{
		mainSchema:  schema,
		joinSchemas: make(map[string]model.Schema),
		builder:     query.NewBuilder(),
		database:    database,
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
	missingCols := q.getMissingCols(columns)
	if len(missingCols) > 0 {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Missing cols: %s", strings.Join(missingCols, ",")))
	}
	q.builder.Selections.AddAll(q.parseCols(columns)...)
	return q
}

func (q *Query) Omit(columns ...string) *Query {
	missingCols := q.getMissingCols(columns)
	if len(missingCols) > 0 {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Missing cols: %s", strings.Join(missingCols, ",")))
	}
	q.builder.Omissions.AddAll(q.parseCols(columns)...)
	return q
}

func (q *Query) join(joinType query.JoinType, otherSchema string, clause query.Clause) {
	schema, err := q.database.GetSchemaByName(otherSchema)
	if err != nil {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Failed to find join schema: %s", otherSchema))
	}
	if _, ok := q.joinSchemas[otherSchema]; !ok {
		q.joinSchemas[otherSchema] = schema
	}
	newJoin := query.Join{Schema: q.mainSchema.Name, OtherSchema: schema.Name, Type: joinType, JoinClause: clause}
	if !newJoin.IsValid(q) {
		q.buildErrors = append(q.buildErrors, fmt.Errorf("Invalid join of type: %s", joinType))
	}
	q.builder.Join(newJoin)
}

func (q *Query) Join(otherSchema string, clause query.Clause) *Query {
	q.join(query.InnerJoin, otherSchema, clause)
	return q
}

func (q *Query) OuterJoin(otherSchema string, clause query.Clause) *Query {
	q.join(query.OuterJoin, otherSchema, clause)
	return q
}

func (q *Query) LeftJoin(otherSchema string, clause query.Clause) *Query {
	q.join(query.LeftJoin, otherSchema, clause)
	return q
}

func (q *Query) RightJoin(otherSchema string, clause query.Clause) *Query {
	q.join(query.RightJoin, otherSchema, clause)
	return q
}

func (q *Query) Where(clause query.Clause) *Query {
	if !clause.IsValid(q) {
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

func (q *Query) OrderBy(orders ...query.Order) *Query {
	for _, order := range orders {
		if order.IsValid(q) {
			q.builder.OrderBy(order)
		} else {
			q.buildErrors = append(q.buildErrors, fmt.Errorf("Invalid order clause"))
		}
	}
	return q
}

func (q *Query) OrderByCol(column string, desc bool) *Query {
	order := query.ColumnOrder{Column: column, Descending: desc}
	return q.OrderBy(order)
}

func (q *Query) OrderByColDistance(column string, target model.SpatialObject, desc bool) *Query {
	order := query.SpatialOrder{Column: column, Target: target, Descending: desc}
	return q.OrderBy(order)
}

func (q *Query) OrderByNearestTo(target model.SpatialObject, desc bool) *Query {
	if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial ordering"))
	} else if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial ordering"))
	}

	if q.mainSchema.LocationFieldNames.Size() == 1 {
		q.OrderBy(query.SpatialOrder{Column: q.mainSchema.LocationFieldNames.Keys()[0], Target: target, Descending: desc})
	} else {
		q.OrderBy(query.SpatialOrder{Column: q.mainSchema.RegionFieldNames.Keys()[0], Target: target, Descending: desc})
	}
	return q
}

func (q *Query) CoveredBy(target model.SpatialObject) *Query {
	if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.mainSchema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.CoveredBy{Column: q.mainSchema.LocationFieldNames.Keys()[0], Target: target})
	} else {
		q.builder.Where(query.CoveredBy{Column: q.mainSchema.RegionFieldNames.Keys()[0], Target: target})
	}
	return q
}

func (q *Query) Covers(target model.SpatialObject) *Query {
	if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.mainSchema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.Covers{Column: q.mainSchema.LocationFieldNames.Keys()[0], Target: target})
	} else {
		q.builder.Where(query.Covers{Column: q.mainSchema.RegionFieldNames.Keys()[0], Target: target})
	}
	return q
}

func (q *Query) WithinRangeOf(targets []model.SpatialObject, rangeMeters float64) *Query {
	if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.mainSchema.LocationFieldNames.Size() == 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else {
		q.builder.Where(query.WithinRangeOf{Column: q.mainSchema.RegionFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
	}
	return q
}

func (q *Query) HasWithinRange(targets []model.SpatialObject, rangeMeters float64) *Query {
	if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() > 1 {
		q.buildErrors = append(q.buildErrors, errors.New("Multiple spatial fields in schema, please specify column for spatial query"))
	} else if q.mainSchema.LocationFieldNames.Size()+q.mainSchema.RegionFieldNames.Size() == 0 {
		q.buildErrors = append(q.buildErrors, errors.New("No spatial fields in schema for spatial query"))
	}

	if q.mainSchema.LocationFieldNames.Size() == 1 {
		q.builder.Where(query.HasWithinRange{Column: q.mainSchema.LocationFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
	} else {
		q.builder.Where(query.HasWithinRange{Column: q.mainSchema.RegionFieldNames.Keys()[0], Targets: targets, Range: rangeMeters})
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
	statement, args := query.CompileSQL(*q.builder, q)
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
	statement, args := query.CompileSQL(*q.builder, q)
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
	statement, args := query.CompileSQL(*q.builder, q)
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
	vals, err := model.ParseObject(object, q.mainSchema)
	if err != nil {
		return nil, err
	}
	q.builder.InsertValues = vals
	statement, args := query.CompileSQL(*q.builder, q)
	return q.database.Execute(statement, args...)
}

func (q *Query) Update(object interface{}) (sql.Result, error) {
	if q.Error() != nil {
		return nil, q.Error()
	}
	q.builder.QueryType = query.UpdateQuery
	vals, err := model.ParseObject(object, q.mainSchema)
	if err != nil {
		return nil, err
	}
	if len(vals) > 1 {
		return nil, errors.New("Can only update 1 record each time")
	}
	q.builder.InsertValues = vals
	statement, args := query.CompileSQL(*q.builder, q)
	return q.database.Execute(statement, args...)
}

func (q *Query) parseCols(columns []string) []string {
	parsedCols := make([]string, 0)
	for _, col := range columns {
		parsedCols = append(parsedCols, q.GetField(col).GetFullName())
	}
	return parsedCols
}

func (q *Query) getMissingCols(columns []string) []string {
	missingCols := make([]string, 0)
	for _, col := range columns {
		if !q.HasField(col) {
			missingCols = append(missingCols, col)
		}
	}
	return missingCols
}

/* Functions for checking/compilation of query */
func (q *Query) splitFieldName(field string) (schema string, fieldName string) {
	vals := strings.Split(field, ".")
	if len(vals) > 2 || len(vals) <= 0 {
		return
	}
	if len(vals) == 1 {
		schema = q.mainSchema.Name
		fieldName = vals[0]
	} else {
		schema = vals[0]
		fieldName = vals[1]
	}
	return
}

func (q *Query) isMainSchema(schema string) bool {
	return schema == q.mainSchema.Name
}

func (q *Query) HasSchema(schema string) bool {
	_, inJoin := q.joinSchemas[schema]
	return schema == q.mainSchema.Name || inJoin
}

func (q *Query) HasField(field string) bool {
	schema, fieldName := q.splitFieldName(field)
	if !q.HasSchema(schema) {
		return false
	}
	if q.isMainSchema(schema) {
		_, ok := q.mainSchema.FieldsByName[fieldName]
		return ok
	}
	_, ok := q.joinSchemas[schema].FieldsByName[fieldName]
	return ok
}

func (q *Query) GetField(field string) *model.Field {
	if !q.HasField(field) {
		return nil
	}
	schema, fieldName := q.splitFieldName(field)
	if q.isMainSchema(schema) {
		return q.mainSchema.FieldsByName[fieldName]
	}
	return q.joinSchemas[schema].FieldsByName[fieldName]
}

func (q *Query) HasFieldOfType(field string, datatype model.DataType) bool {
	if !q.HasField(field) {
		return false
	}
	return q.GetField(field).DataType == datatype
}

func (q *Query) GetMainSchema() model.Schema {
	return q.mainSchema
}

func (q *Query) GetJoinSchemas() map[string]model.Schema {
	return q.joinSchemas
}

func (q *Query) GetAdapterInfo() adapter.AdapterInfo {
	return q.database.adapter
}
