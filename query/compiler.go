package query

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/utils"
)

type Compiler struct {
	AdapterInfo adapter.AdapterInfo
	Schema      model.Schema
}

func (c Compiler) parseSelectionField(name string) string {
	field := c.Schema.FieldsByName[name]
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if c.AdapterInfo.SpatialType() == adapter.PostGisExtension {
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
		if c.AdapterInfo.SpatialType() == adapter.PostGisExtension {
			return "ST_GeomFromGeoJSON(?)::geography"
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
	var targetFieldsSet *utils.Set
	if builder.Selections.Size() == 0 {
		targetSet := utils.NewSet()
		for _, field := range c.Schema.Fields {
			targetSet.Add(field.Name)
		}
		targetFieldsSet = targetSet.Difference(builder.Omissions)
	} else {
		targetFieldsSet = builder.Selections.Difference(builder.Omissions)
	}

	if builder.QueryType == UpdateQuery && len(builder.Clauses) == 0 {
		targetFieldsSet = targetFieldsSet.Difference(c.Schema.PrimaryFieldNames)
	}

	targetFields := targetFieldsSet.Keys()

	switch qType := builder.QueryType; qType {
	case SelectQuery:
		sql.WriteString("SELECT ")

		if builder.IsCount {
			sql.WriteString("COUNT(*) ")
		} else {
			selBuilder := strings.Builder{}
			for i, sel := range targetFields {
				selBuilder.WriteString(c.parseSelectionField(sel))
				if i < len(targetFields)-1 {
					selBuilder.WriteString(",")
				}
			}
			selBuilder.WriteString(" ")
			sql.WriteString(selBuilder.String())
		}

		sql.WriteString("FROM ")

	case InsertQuery:
		sql.WriteString("INSERT INTO ")

	case UpdateQuery:
		sql.WriteString("UPDATE ")
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
			clauseSql, clauseValues := clause.Sql(c.Schema.FieldsByName, c.AdapterInfo.SpatialType())
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
		for i, key := range targetFields {
			sql.WriteString(c.Schema.FieldsByName[key].DBName)
			if i < len(targetFields)-1 {
				sql.WriteString(",")
			}
		}
		sql.WriteString(") ")
		sql.WriteString("VALUES ")

		for i, insertVal := range builder.InsertValues {
			sql.WriteString("(")
			for k, key := range targetFields {
				sql.WriteString(c.parseInsertionValuePlaceholder(key))
				values = append(values, insertVal[key])
				if k < len(targetFields)-1 {
					sql.WriteString(",")
				}
			}
			sql.WriteString(")")
			if i < len(builder.InsertValues)-1 {
				sql.WriteString(",")
			}
		}
	case UpdateQuery:
		insertVal := builder.InsertValues[0]
		sql.WriteString("SET ")
		for i, key := range targetFields {
			sql.WriteString(c.Schema.FieldsByName[key].DBName)
			sql.WriteString(" = ")
			sql.WriteString(c.parseInsertionValuePlaceholder(key))
			values = append(values, insertVal[key])
			if i < len(targetFields)-1 {
				sql.WriteString(",")
			}
		}
		sql.WriteString(" WHERE ")
		if len(builder.Clauses) > 0 {
			clause := builder.Clauses[0]
			if len(builder.Clauses) > 1 {
				clause = append(And{}, builder.Clauses...)
			}
			clauseSql, clauseValues := clause.Sql(c.Schema.FieldsByName, c.AdapterInfo.SpatialType())
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		} else {
			primaryClauses := And{}
			for _, field := range c.Schema.PrimaryFields {
				primaryClauses = append(primaryClauses, Equal{Column: field.Name, Value: insertVal[field.Name]})
			}
			clauseSql, clauseValues := primaryClauses.Sql(c.Schema.FieldsByName, c.AdapterInfo.SpatialType())
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		}
	}
	sql.WriteString(";")
	return replacePlaceholder(sql.String(), c.AdapterInfo.Placeholder()), values
}
func replacePlaceholder(sqlString string, style adapter.PlaceholderStyle) string {
	switch style {
	case adapter.DollarPlaceholder:
		for nParam := 1; strings.Contains(sqlString, "?"); nParam++ {
			sqlString = strings.Replace(sqlString, "?", fmt.Sprintf("$%d", nParam), 1)
		}
	default:
		return sqlString
	}
	return sqlString
}
