package query

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
)

type Clause interface {
	Condition() string
	IsValid(info QueryInfo) bool
	Sql(info QueryInfo) (string, []interface{})
}

type GreaterThan struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e GreaterThan) Sql(info QueryInfo) (string, []interface{}) {
	if e.OtherColumn != "" {
		return fmt.Sprintf("%s > %s", info.GetField(e.Column).GetFullDBName(), info.GetField(e.OtherColumn).GetFullDBName()), []interface{}{}
	}
	return fmt.Sprintf("%s > ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e GreaterThan) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType != model.LocationType && field.DataType != model.RegionType
	secondOk := true
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType != model.LocationType && otherField.DataType != model.RegionType
	}
	return firstOk && secondOk
}

func (e GreaterThan) Condition() string {
	return ">"
}

type LessThan struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e LessThan) Sql(info QueryInfo) (string, []interface{}) {
	if e.OtherColumn != "" {
		return fmt.Sprintf("%s < %s", info.GetField(e.Column).GetFullDBName(), info.GetField(e.OtherColumn).GetFullDBName()), []interface{}{}
	}
	return fmt.Sprintf("%s < ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e LessThan) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType != model.LocationType && field.DataType != model.RegionType
	secondOk := true
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType != model.LocationType && otherField.DataType != model.RegionType
	}
	return firstOk && secondOk
}

func (e LessThan) Condition() string {
	return "<"
}

type Equal struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e Equal) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	field := info.GetField(e.Column)
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if spatialType == adapter.PostGisExtension {
			if e.OtherColumn != "" {
				otherField := info.GetField(e.OtherColumn)
				return fmt.Sprintf("ST_Equals(%s::geometry, %s::geometry)", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
			}
			return fmt.Sprintf("ST_Equals(%s::geometry, ST_GeomFromGeoJSON(?))", field.GetFullDBName()), []interface{}{e.Value}
		} else {
			if e.OtherColumn != "" {
				otherField := info.GetField(e.OtherColumn)
				return fmt.Sprintf("%s = %s", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
			}
			return fmt.Sprintf("%s = ?", field.GetFullDBName()), []interface{}{e.Value}
		}
	default:
		if e.OtherColumn != "" {
			otherField := info.GetField(e.OtherColumn)
			return fmt.Sprintf("%s = %s", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
		}
		return fmt.Sprintf("%s = ?", field.GetFullDBName()), []interface{}{e.Value}
	}
}

func (e Equal) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
	}
	if e.Value != nil {
		switch e.Value.(type) {
		case model.Location:
			return field.DataType == model.LocationType
		case model.Region:
			return field.DataType == model.RegionType
		default:
			return true
		}
	}

	return true
}

func (e Equal) Condition() string {
	return "="
}

type NotEqual struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e NotEqual) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	field := info.GetField(e.Column)
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if spatialType == adapter.PostGisExtension {
			if e.OtherColumn != "" {
				otherField := info.GetField(e.OtherColumn)
				return fmt.Sprintf("NOT ST_Equals(%s::geometry, %s::geometry)", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
			}
			return fmt.Sprintf("NOT ST_Equals(%s::geometry, ST_GeomFromGeoJSON(?))", field.GetFullDBName()), []interface{}{e.Value}
		} else {
			if e.OtherColumn != "" {
				otherField := info.GetField(e.OtherColumn)
				return fmt.Sprintf("%s <> %s", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
			}
			return fmt.Sprintf("%s <> ?", field.GetFullDBName()), []interface{}{e.Value}
		}
	default:
		if e.OtherColumn != "" {
			otherField := info.GetField(e.OtherColumn)
			return fmt.Sprintf("%s <> %s", field.GetFullDBName(), otherField.GetFullDBName()), []interface{}{}
		}
		return fmt.Sprintf("%s <> ?", field.GetFullDBName()), []interface{}{e.Value}
	}
}

func (e NotEqual) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
	}
	if e.Value != nil {
		switch e.Value.(type) {
		case model.Location:
			return field.DataType == model.LocationType
		case model.Region:
			return field.DataType == model.RegionType
		default:
			return true
		}
	}

	return true
}

func (e NotEqual) Condition() string {
	return "<>"
}

type GreaterThanOrEqual struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e GreaterThanOrEqual) Sql(info QueryInfo) (string, []interface{}) {
	if e.OtherColumn != "" {
		return fmt.Sprintf("%s >= %s", info.GetField(e.Column).GetFullDBName(), info.GetField(e.OtherColumn).GetFullDBName()), []interface{}{}
	}
	return fmt.Sprintf("%s >= ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e GreaterThanOrEqual) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType != model.LocationType && field.DataType != model.RegionType
	secondOk := true
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType != model.LocationType && otherField.DataType != model.RegionType
	}
	return firstOk && secondOk
}

func (e GreaterThanOrEqual) Condition() string {
	return ">="
}

type LessThanOrEqual struct {
	Column      string
	OtherColumn string
	Value       interface{}
}

func (e LessThanOrEqual) Sql(info QueryInfo) (string, []interface{}) {
	if e.OtherColumn != "" {
		return fmt.Sprintf("%s <= %s", info.GetField(e.Column).GetFullDBName(), info.GetField(e.OtherColumn).GetFullDBName()), []interface{}{}
	}
	return fmt.Sprintf("%s <= ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e LessThanOrEqual) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType != model.LocationType && field.DataType != model.RegionType
	secondOk := true
	if e.OtherColumn != "" {
		otherField := info.GetField(e.OtherColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType != model.LocationType && otherField.DataType != model.RegionType
	}
	return firstOk && secondOk
}

func (e LessThanOrEqual) Condition() string {
	return "<="
}

type Like struct {
	Column string
	Value  string
}

func (e Like) Sql(info QueryInfo) (string, []interface{}) {
	return fmt.Sprintf("%s LIKE ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e Like) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	} else {
		return field.DataType == model.String
	}
}

func (e Like) Condition() string {
	return "LIKE"
}

type NotLike struct {
	Column string
	Value  string
}

func (e NotLike) Sql(info QueryInfo) (string, []interface{}) {
	return fmt.Sprintf("%s NOT LIKE ?", info.GetField(e.Column).GetFullDBName()), []interface{}{e.Value}
}

func (e NotLike) IsValid(info QueryInfo) bool {
	field := info.GetField(e.Column)
	if field == nil {
		return false
	} else {
		return field.DataType == model.String
	}
}

func (e NotLike) Condition() string {
	return "NOT LIKE"
}

// Geography specific clauses
type CoveredBy struct {
	Column       string
	TargetColumn string
	Target       model.SpatialObject
}

func (c CoveredBy) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	if spatialType == adapter.PostGisExtension {
		if c.TargetColumn != "" {
			return fmt.Sprintf("ST_Covers(%s, %s)", info.GetField(c.TargetColumn).GetFullDBName(), info.GetField(c.Column).GetFullDBName()), []interface{}{}
		}
		return fmt.Sprintf("ST_Covers(ST_GeomFromGeoJSON(?)::geography, %s)", info.GetField(c.Column).GetFullDBName()), []interface{}{c.Target}
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (c CoveredBy) IsValid(info QueryInfo) bool {
	field := info.GetField(c.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType == model.LocationType || field.DataType == model.RegionType
	secondOk := true
	if c.TargetColumn != "" {
		otherField := info.GetField(c.TargetColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType == model.LocationType || otherField.DataType == model.RegionType
	}
	return firstOk && secondOk
}

func (c CoveredBy) Condition() string {
	return "CoveredBy"
}

type Covers struct {
	Column       string
	TargetColumn string
	Target       model.SpatialObject
}

func (c Covers) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	if spatialType == adapter.PostGisExtension {
		if c.TargetColumn != "" {
			return fmt.Sprintf("ST_Covers(%s, %s)", info.GetField(c.Column).GetFullDBName(), info.GetField(c.TargetColumn).GetFullDBName()), []interface{}{}
		}
		return fmt.Sprintf("ST_Covers(%s, ST_GeomFromGeoJSON(?)::geography)", info.GetField(c.Column).GetFullDBName()), []interface{}{c.Target}
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (c Covers) IsValid(info QueryInfo) bool {
	field := info.GetField(c.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType == model.LocationType || field.DataType == model.RegionType
	secondOk := true
	if c.TargetColumn != "" {
		otherField := info.GetField(c.TargetColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType == model.LocationType || otherField.DataType == model.RegionType
	}
	return firstOk && secondOk
}

func (c Covers) Condition() string {
	return "Covers"
}

type WithinRangeOf struct {
	Column       string
	TargetColumn string
	Targets      []model.SpatialObject
	Range        float64
}

func (w WithinRangeOf) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	if spatialType == adapter.PostGisExtension {
		sql := strings.Builder{}
		vals := []interface{}{}

		if w.TargetColumn != "" {
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, %s, ?)", info.GetField(w.Column).GetFullDBName(), info.GetField(w.TargetColumn).GetFullDBName()))
			vals = append(vals, w.Range)
		}

		for _, targetObj := range w.Targets {
			if sql.Len() > 0 {
				sql.WriteString(" OR ")
			}
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, ST_GeomFromGeoJSON(?)::geography, ?)", info.GetField(w.Column).GetFullDBName()))
			vals = append(vals, targetObj, w.Range)
		}
		return sql.String(), vals
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (w WithinRangeOf) IsValid(info QueryInfo) bool {
	field := info.GetField(w.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType == model.LocationType || field.DataType == model.RegionType
	secondOk := true
	if w.TargetColumn != "" {
		otherField := info.GetField(w.TargetColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType == model.LocationType || otherField.DataType == model.RegionType
	}
	return firstOk && secondOk
}

func (w WithinRangeOf) Condition() string {
	return "WithinRangeOf"
}

type HasWithinRange struct {
	Column       string
	TargetColumn string
	Targets      []model.SpatialObject
	Range        float64
}

func (h HasWithinRange) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	if spatialType == adapter.PostGisExtension {
		sql := strings.Builder{}
		vals := []interface{}{}

		if h.TargetColumn != "" {
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, %s, ?)", info.GetField(h.Column).GetFullDBName(), info.GetField(h.TargetColumn).GetFullDBName()))
			vals = append(vals, h.Range)
		}

		for _, targetObj := range h.Targets {
			if sql.Len() > 0 {
				sql.WriteString(" AND ")
			}
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s, ST_GeomFromGeoJSON(?)::geography, ?)", info.GetField(h.Column).GetFullDBName()))
			vals = append(vals, targetObj, h.Range)
		}
		return sql.String(), vals
	} else {
		// Not implemented
		return "", []interface{}{}
	}
}

func (h HasWithinRange) IsValid(info QueryInfo) bool {
	field := info.GetField(h.Column)
	if field == nil {
		return false
	}
	firstOk := field.DataType == model.LocationType || field.DataType == model.RegionType
	secondOk := true
	if h.TargetColumn != "" {
		otherField := info.GetField(h.TargetColumn)
		if otherField == nil {
			return false
		}
		secondOk = otherField.DataType == model.LocationType || otherField.DataType == model.RegionType
	}
	return firstOk && secondOk
}

func (h HasWithinRange) Condition() string {
	return "HasWithinRange"
}

type Or []Clause

func (e Or) Sql(info QueryInfo) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql(info)
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

func (e Or) IsValid(info QueryInfo) bool {
	for _, clause := range e {
		if !clause.IsValid(info) {
			return false
		}
	}
	return true
}

func (e Or) Condition() string {
	return "OR"
}

type And []Clause

func (e And) Sql(info QueryInfo) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql(info)
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

func (e And) IsValid(info QueryInfo) bool {
	for _, clause := range e {
		if !clause.IsValid(info) {
			return false
		}
	}
	return true
}

func (e And) Condition() string {
	return "AND"
}
