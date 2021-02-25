package query

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
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
			return fmt.Sprintf("ST_AsGeoJSON(%s)", field.DBName)
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
	switch qType := builder.QueryType; qType {
	case SelectQuery:
		sql.WriteString("SELECT ")

		if builder.IsCount {
			sql.WriteString("COUNT(*) ")
		} else {
			selection := "* "
			if len(builder.Selections) > 0 {
				selBuilder := strings.Builder{}
				for i, sel := range builder.Selections {
					selBuilder.WriteString(c.parseSelectionField(sel))
					if i < len(builder.Selections)-1 {
						selBuilder.WriteString(",")
					}
				}
				selBuilder.WriteString(" ")
				selection = selBuilder.String()
			}
			// TODO: Insert selection fields here
			sql.WriteString(selection)
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
		if len(builder.Selections) > 0 {
			sql.WriteString("(")
			for i, key := range builder.Selections {
				sql.WriteString(c.Schema.FieldsByName[key].DBName)
				if i < len(builder.Selections)-1 {
					sql.WriteString(",")
				}
			}
			sql.WriteString(") ")
		}
		sql.WriteString("VALUES ")

		if len(builder.Selections) > 0 {
			for i, insertVal := range builder.InsertValues {
				sql.WriteString("(")
				for k, key := range builder.Selections {
					sql.WriteString(c.parseInsertionValuePlaceholder(key))
					values = append(values, insertVal[key])
					if k < len(builder.Selections)-1 {
						sql.WriteString(",")
					}
				}
				sql.WriteString(")")
				if i < len(builder.InsertValues)-1 {
					sql.WriteString(",")
				}
			}
		} else {
			for i, insertVal := range builder.InsertValues {
				sql.WriteString("(")
				for k, field := range c.Schema.Fields {
					sql.WriteString(c.parseInsertionValuePlaceholder(field.Name))
					values = append(values, insertVal[field.Name])
					if k < len(c.Schema.Fields)-1 {
						sql.WriteString(",")
					}
				}
				sql.WriteString(")")
				if i < len(builder.InsertValues)-1 {
					sql.WriteString(",")
				}
			}
		}
	}
	sql.WriteString(";")
	return sql.String(), values
}
