package query

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/utils"
)

type Compiler struct {
	info QueryInfo
}

func CompileSQL(builder Builder, info QueryInfo) (string, []interface{}) {
	compiler := Compiler{info: info}
	return compiler.compileSQL(builder)
}

func CompileTableCreation(info QueryInfo, ifNotExists bool) string {
	compiler := Compiler{info: info}
	return compiler.compileTableCreation(ifNotExists)
}

func CompileIndexCreation(info QueryInfo) []string {
	compiler := Compiler{info: info}
	return compiler.compileIndexCreation()
}

func (c Compiler) parseSelectionField(name string) string {
	field := c.info.GetField(name)
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if c.info.GetAdapterInfo().SpatialType() == adapter.PostGisExtension {
			return fmt.Sprintf("ST_AsGeoJSON(%s) as %s", field.GetFullDBName(), field.DBName)
		} else {
			return field.GetFullDBName()
		}
	default:
		return field.GetFullDBName()
	}
}

func (c Compiler) parseInsertionValuePlaceholder(name string) string {
	field := c.info.GetField(name)
	switch field.DataType {
	case model.LocationType, model.RegionType:
		if c.info.GetAdapterInfo().SpatialType() == adapter.PostGisExtension {
			return "ST_GeomFromGeoJSON(?)::geography"
		} else {
			return "?"
		}
	default:
		return "?"
	}
}

func (c Compiler) compileSQL(builder Builder) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	var targetFieldsSet *utils.Set
	if builder.Selections.Size() == 0 {
		targetSet := c.info.GetMainSchema().AllFieldNames
		for _, schema := range c.info.GetJoinSchemas() {
			targetSet = targetSet.Union(schema.AllFieldNames)
		}
		targetFieldsSet = targetSet.Difference(builder.Omissions)
	} else {
		targetFieldsSet = builder.Selections.Difference(builder.Omissions)
	}

	if builder.QueryType == UpdateQuery && len(builder.Clauses) == 0 {
		targetFieldsSet = targetFieldsSet.Difference(c.info.GetMainSchema().PrimaryFieldNames)
	}

	targetFields := targetFieldsSet.Keys()

	switch qType := builder.QueryType; qType {
	case SelectQuery:
		sql.WriteString("SELECT ")

		if builder.IsCount {
			if builder.IsDistinct {
				sql.WriteString("COUNT(DISTINCT(")
				selBuilder := strings.Builder{}
				for i, sel := range targetFields {
					selBuilder.WriteString(c.parseSelectionField(sel))
					if i < len(targetFields)-1 {
						selBuilder.WriteString(",")
					}
				}
				sql.WriteString(selBuilder.String())
				sql.WriteString(")) ")
			} else {
				sql.WriteString("COUNT(*) ")
			}
		} else {
			if builder.IsDistinct {
				sql.WriteString("DISTINCT ")
			}
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

	sql.WriteString(c.info.GetMainSchema().Table + " ")

	switch qType := builder.QueryType; qType {
	case SelectQuery:
		for _, join := range builder.Joins {
			joinSql, joinVals := join.Sql(c.info)
			sql.WriteString(joinSql)
			values = append(values, joinVals...)
			sql.WriteString(" ")
		}
		if len(builder.Clauses) > 0 {
			sql.WriteString("WHERE ")
			clause := builder.Clauses[0]
			if len(builder.Clauses) > 1 {
				clause = append(And{}, builder.Clauses...)
			}
			clauseSql, clauseValues := clause.Sql(c.info)
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		}
		if len(builder.Orders) > 0 {
			sql.WriteString(" ORDER BY ")
			for i, order := range builder.Orders {
				clauseSql, clauseValues := order.Sql(c.info)
				sql.WriteString(clauseSql)
				values = append(values, clauseValues...)
				if i < len(builder.Orders)-1 {
					sql.WriteString(",")
				}
			}
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
			sql.WriteString(c.info.GetField(key).DBName)
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
				values = append(values, insertVal[c.info.GetField(key).Name])
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
			sql.WriteString(c.info.GetField(key).DBName)
			sql.WriteString(" = ")
			sql.WriteString(c.parseInsertionValuePlaceholder(key))
			values = append(values, insertVal[c.info.GetField(key).Name])
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
			clauseSql, clauseValues := clause.Sql(c.info)
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		} else {
			primaryClauses := And{}
			for _, field := range c.info.GetMainSchema().PrimaryFields {
				primaryClauses = append(primaryClauses, Equal{Column: field.Name, Value: insertVal[field.Name]})
			}
			clauseSql, clauseValues := primaryClauses.Sql(c.info)
			sql.WriteString(clauseSql)
			values = append(values, clauseValues...)
		}
	}
	sql.WriteString(";")
	return replacePlaceholder(sql.String(), c.info.GetAdapterInfo().Placeholder()), values
}

func (c Compiler) compileTableCreation(ifNotExists bool) string {
	schema := c.info.GetMainSchema()
	sql := strings.Builder{}
	sql.WriteString("CREATE TABLE ")
	if ifNotExists {
		sql.WriteString("IF NOT EXISTS ")
	}
	sql.WriteString(schema.Table)
	sql.WriteString(" (")

	for i, field := range schema.Fields {
		sql.WriteString(field.DBName)
		sql.WriteString(" ")
		sql.WriteString(c.parseFieldType(field))
		qualifiers := c.parseFieldQualifiers(field)
		sql.WriteString(qualifiers)
		if i < len(schema.Fields)-1 {
			sql.WriteString(", ")
		}
	}

	sql.WriteString(");")

	return replacePlaceholder(sql.String(), c.info.GetAdapterInfo().Placeholder())
}

func (c Compiler) compileIndexCreation() []string {
	schema := c.info.GetMainSchema()
	allStatements := make([]string, 0)

	// Create indexes for spatial types
	for _, fieldName := range schema.LocationFieldNames.Keys() {
		field := c.info.GetField(fieldName)
		indexName := fmt.Sprintf("idx_%s_%s", schema.Table, field.DBName)
		allStatements = append(allStatements, fmt.Sprintf("CREATE INDEX %s ON %s USING GIST (%s);", indexName, schema.Table, field.DBName))
	}

	for _, fieldName := range schema.RegionFieldNames.Keys() {
		field := c.info.GetField(fieldName)
		indexName := fmt.Sprintf("idx_%s_%s", schema.Table, field.DBName)
		allStatements = append(allStatements, fmt.Sprintf("CREATE INDEX %s ON %s USING GIST (%s);", indexName, schema.Table, field.DBName))
	}

	return allStatements
}

func (c Compiler) parseFieldQualifiers(field *model.Field) string {
	qualifiers := strings.Builder{}
	adapterInfo := c.info.GetAdapterInfo()
	switch adapterInfo.DatabaseType() {
	case adapter.PostgreSQL:
		if field.PrimaryKey {
			qualifiers.WriteString(" PRIMARY KEY")
		} else {
			if field.NotNull {
				qualifiers.WriteString(" NOT NULL")
			}
			if field.Unique {
				qualifiers.WriteString(" UNIQUE")
			}
		}
		if !field.AutoIncrement && field.HasDefaultValue {
			if field.DefaultValue == nil {
				qualifiers.WriteString(" DEFAULT NULL")
			} else {
				switch v := field.DefaultValue.(type) {
				case model.Location, model.Region:
					val, err := v.(driver.Valuer).Value()
					if err == nil {
						switch adapterInfo.SpatialType() {
						case adapter.PostGisExtension:
							geom := bytes.NewBuffer(val.([]uint8))
							qualifiers.WriteString(fmt.Sprintf(" DEFAULT ST_GeomFromGeoJSON('%v')::geography", geom.String()))
						default:
							qualifiers.WriteString(fmt.Sprintf(" DEFAULT %v", val))
						}
					}
				case driver.Valuer:
					val, err := v.Value()
					if err == nil {
						qualifiers.WriteString(fmt.Sprintf(" DEFAULT %v", val))
					}
				default:
					qualifiers.WriteString(fmt.Sprintf(" DEFAULT %v", toString(v)))
				}
			}
		}
	}
	return qualifiers.String()
}

func (c Compiler) parseFieldType(field *model.Field) string {
	switch c.info.GetAdapterInfo().DatabaseType() {
	case adapter.PostgreSQL:
		if field.AutoIncrement {
			return "serial"
		}
		switch field.DataType {
		case model.Bool:
			return "bool"
		case model.Int:
			return "int"
		case model.Uint:
			return "int" // unsigned not supported by postgresql
		case model.String:
			return "varchar(255)" // TODO: might not be long enough, allow spec
		case model.Float:
			return "float64"
		case model.Time:
			return "time"
		case model.Bytes:
			return "bytea"
		case model.LocationType:
			return "geography(point)"
		case model.RegionType:
			return "geography(polygon)"
		case model.TimestampType:
			return "timestamp"
		default:
			return string(field.DataType)
		}
	default:
		return ""
	}
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

func toString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	}
	return ""
}
