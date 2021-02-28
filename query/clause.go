package query

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
)

type Clause interface {
	Condition() string
	Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{})
}

type GreaterThan struct {
	Column string
	Value  string
}

func (e GreaterThan) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s > ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e GreaterThan) Condition() string {
	return ">"
}

type LessThan struct {
	Column string
	Value  string
}

func (e LessThan) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s < ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e LessThan) Condition() string {
	return "<"
}

type Equal struct {
	Column string
	Value  interface{}
}

func (e Equal) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	switch fields[e.Column].DataType {
	case model.LocationType, model.RegionType:
		if spatialType == adapter.PostGisExtension {
			return fmt.Sprintf("ST_Equals(%s::geometry, ST_GeomFromGeoJSON(?))", fields[e.Column].DBName), []interface{}{e.Value}
		} else {
			return fmt.Sprintf("%s = ?", fields[e.Column].DBName), []interface{}{e.Value}
		}
	default:
		return fmt.Sprintf("%s = ?", fields[e.Column].DBName), []interface{}{e.Value}
	}
}

func (e Equal) Condition() string {
	return "="
}

type GreaterThanOrEqual struct {
	Column string
	Value  string
}

func (e GreaterThanOrEqual) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s >= ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e GreaterThanOrEqual) Condition() string {
	return ">="
}

type LessThanOrEqual struct {
	Column string
	Value  string
}

func (e LessThanOrEqual) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s <= ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e LessThanOrEqual) Condition() string {
	return "<="
}

type NotEqual struct {
	Column string
	Value  string
}

func (e NotEqual) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	switch fields[e.Column].DataType {
	case model.LocationType, model.RegionType:
		if spatialType == adapter.PostGisExtension {
			return fmt.Sprintf("NOT ST_Equals(%s::geometry, ST_GeomFromGeoJSON(?))", fields[e.Column].DBName), []interface{}{e.Value}
		} else {
			return fmt.Sprintf("%s <> ?", fields[e.Column].DBName), []interface{}{e.Value}
		}
	default:
		return fmt.Sprintf("%s <> ?", fields[e.Column].DBName), []interface{}{e.Value}
	}
}

func (e NotEqual) Condition() string {
	return "<>"
}

type Like struct {
	Column string
	Value  string
}

func (e Like) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s LIKE ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e Like) Condition() string {
	return "LIKE"
}

type NotLike struct {
	Column string
	Value  string
}

func (e NotLike) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	return fmt.Sprintf("%s NOT LIKE ?", fields[e.Column].DBName), []interface{}{e.Value}
}

func (e NotLike) Condition() string {
	return "NOT LIKE"
}

// Geography specific clauses
type CoveredBy struct {
	Column string
	Target model.SpatialObject
}

func (c CoveredBy) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	if spatialType == adapter.PostGisExtension {
		return fmt.Sprintf("ST_Covers(ST_GeomFromGeoJSON(?)::geography, %s)", fields[c.Column].DBName), []interface{}{c.Target}
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (c CoveredBy) Condition() string {
	return "CoveredBy"
}

type Covers struct {
	Column string
	Target model.SpatialObject
}

func (c Covers) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	if spatialType == adapter.PostGisExtension {
		return fmt.Sprintf("ST_Covers(%s, ST_GeomFromGeoJSON(?)::geography)", fields[c.Column].DBName), []interface{}{c.Target}
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (c Covers) Condition() string {
	return "Covers"
}

type WithinRangeOf struct {
	Column  string
	Targets []model.SpatialObject
	Range   float64
}

func (w WithinRangeOf) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	if spatialType == adapter.PostGisExtension {
		sql := strings.Builder{}
		vals := []interface{}{}
		for i, targetObj := range w.Targets {
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, ST_GeomFromGeoJSON(?)::geography, ?)", fields[w.Column].DBName))
			vals = append(vals, targetObj, w.Range)
			if i < len(w.Targets)-1 {
				sql.WriteString(" OR ")
			}
		}
		return sql.String(), vals
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (w WithinRangeOf) Condition() string {
	return "WithinRangeOf"
}

type HasWithinRange struct {
	Column  string
	Targets []model.SpatialObject
	Range   float64
}

func (h HasWithinRange) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	if spatialType == adapter.PostGisExtension {
		sql := strings.Builder{}
		vals := []interface{}{}
		for i, targetObj := range h.Targets {
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, ST_GeomFromGeoJSON(?)::geography, ?)", fields[h.Column].DBName))
			vals = append(vals, targetObj, h.Range)
			if i < len(h.Targets)-1 {
				sql.WriteString(" AND ")
			}
		}
		return sql.String(), vals
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (h HasWithinRange) Condition() string {
	return "HasWithinRange"
}

type Or []Clause

func (e Or) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql(fields, spatialType)
		if i == 0 {
			sql.WriteString(clauseSql)
		} else {
			sql.WriteString(" OR ")
			sql.WriteString(clauseSql)
		}
		values = append(values, clauseVals...)
	}
	return sql.String(), values
}

func (e Or) Condition() string {
	return "OR"
}

type And []Clause

func (e And) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql(fields, spatialType)
		if i == 0 {
			sql.WriteString(clauseSql)
		} else {
			sql.WriteString(" AND ")
			sql.WriteString(clauseSql)
		}
		values = append(values, clauseVals...)
	}
	return sql.String(), values
}

func (e And) Condition() string {
	return "AND"
}
