package atlas

import (
	"fmt"
	"strings"

	"github.com/JayPeeTeeDee/atlas/query"
)

func CompileSQL(builder query.Builder) (string, []interface{}) {
	if builder.QueryType == query.SelectQuery {
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
				sql.WriteString(clauseSql + " ")
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
			// TODO: implement
			// if len(builder.Clauses) > 0 {
			// 	sql.WriteString("WHERE ")
			// 	clause := builder.Clauses[0]
			// 	if len(builder.Clauses) > 1 {
			// 		clause = append(query.And{}, builder.Clauses...)
			// 	}
			// 	clauseSql, clauseValues := clause.Sql()
			// 	sql.WriteString(clauseSql + " ")
			// 	values = append(values, clauseValues...)
			// 	}
			// }
		}
		return sql.String(), values
	}
	return "", nil
}
