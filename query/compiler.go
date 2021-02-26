package query

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/utils"
)

type Compiler struct {
	SpatialType adapter.SpatialExtension
	Schema      model.Schema
}

func (c Compiler) parseSelectionField(name string) string {
	field := c.Schema.FieldsByName[name]
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if c.SpatialType == adapter.PostGisExtension {
			return fmt.Sprintf("ST_AsGeoJSON(%s) as %s", field.DBName, field.DBName)
		} else {
			return field.DBName
		}
	default:
		return field.DBName
	}
}

func (c Compiler) parseInsertionValuePlaceholder(name string) string {
	field := c.Schema.FieldsByName[name]
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if c.SpatialType == adapter.PostGisExtension {
			return "ST_GeomFromGeoJSON(?)"
		} else {
			return "?"
		}
	default:
		return "?"
	}
}

func (c Compiler) CompileSQL(builder Builder) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	var target_fields []string
	if builder.Selections.Size() == 0 {
		target_set := utils.NewSet()
		for _, field := range c.Schema.Fields {
			target_set.Add(field.Name)
		}
		target_fields = target_set.Difference(builder.Omissions).Keys()
	} else {
		target_fields = builder.Selections.Difference(builder.Omissions).Keys()
	}

	switch qType := builder.QueryType; qType {
	case SelectQuery:
		sql.WriteString("SELECT ")

		if builder.IsCount {
			sql.WriteString("COUNT(*) ")
		} else {
			selBuilder := strings.Builder{}
			for i, sel := range target_fields {
				selBuilder.WriteString(c.parseSelectionField(sel))
				if i < len(target_fields)-1 {
					selBuilder.WriteString(",")
				}
			}
			selBuilder.WriteString(" ")
			sql.WriteString(selBuilder.String())
		}

		sql.WriteString("FROM ")

	case InsertQuery:
		sql.WriteString("INSERT INTO ")
	}

	sql.WriteString(c.Schema.Table + " ")

	switch qType := builder.QueryType; qType {
	case SelectQuery:
		if len(builder.Clauses) > 0 {
			sql.WriteString("WHERE ")
			clause := builder.Clauses[0]
			if len(builder.Clauses) > 1 {
				clause = append(And{}, builder.Clauses...)
			}
			clauseSql, clauseValues := clause.Sql(c.Schema.FieldsByName, c.SpatialType)
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		}
		if builder.Limit > 0 {
			sql.WriteString(" LIMIT ")
			sql.WriteString(fmt.Sprintf("%d", builder.Limit))
		}
		if builder.Offset > 0 {
			sql.WriteString(" OFFSET ")
			sql.WriteString(fmt.Sprintf("%d", builder.Offset))
		}

	case InsertQuery:
		sql.WriteString("(")
		for i, key := range target_fields {
			sql.WriteString(c.Schema.FieldsByName[key].DBName)
			if i < len(target_fields)-1 {
				sql.WriteString(",")
			}
		}
		sql.WriteString(") ")
		sql.WriteString("VALUES ")

		for i, insertVal := range builder.InsertValues {
			sql.WriteString("(")
			for k, key := range target_fields {
				sql.WriteString(c.parseInsertionValuePlaceholder(key))
				values = append(values, insertVal[key])
				if k < len(target_fields)-1 {
					sql.WriteString(",")
				}
			}
			sql.WriteString(")")
			if i < len(builder.InsertValues)-1 {
				sql.WriteString(",")
			}
		}
	}
	sql.WriteString(";")
	return sql.String(), values
}
