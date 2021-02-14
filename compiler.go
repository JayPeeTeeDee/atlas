package atlas

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/query"
)

func CompileSQL(builder query.Builder) (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	switch qType := builder.QueryType; qType {
	case query.SelectQuery:
		sql.WriteString("SELECT ")

		if builder.IsCount {
			sql.WriteString("COUNT(*) ")
		} else {
			selection := "* "
			// TODO: Insert selection fields here
			sql.WriteString(selection)
		}

		sql.WriteString("FROM ")

	case query.InsertQuery:
		sql.WriteString("INSERT INTO ")
	}

	sql.WriteString(builder.TableName + " ")

	switch qType := builder.QueryType; qType {
	case query.SelectQuery:
		if len(builder.Clauses) > 0 {
			sql.WriteString("WHERE ")
			clause := builder.Clauses[0]
			if len(builder.Clauses) > 1 {
				clause = append(query.And{}, builder.Clauses...)
			}
			clauseSql, clauseValues := clause.Sql()
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

	case query.InsertQuery:
		if len(builder.Selections) > 0 {
			sql.WriteString("(")
			sql.WriteString(strings.Join(builder.Selections, ","))
			sql.WriteString(") ")
		}
		sql.WriteString("VALUES ")

		if len(builder.Selections) > 0 {
			for i, insertVal := range builder.InsertValues {
				sql.WriteString("(")
				for k, key := range builder.Selections {
					sql.WriteString("?")
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
			for i, insertAllVals := range builder.FieldValues {
				sql.WriteString("(")
				for k, val := range insertAllVals {
					sql.WriteString("?")
					values = append(values, val)
					if k < len(insertAllVals)-1 {
						sql.WriteString(",")
					}
				}
				sql.WriteString(")")
				if i < len(builder.FieldValues)-1 {
					sql.WriteString(",")
				}
			}
		}
	}
	sql.WriteString(";")
	return sql.String(), values
}
