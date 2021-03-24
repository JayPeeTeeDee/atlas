package query

import "fmt"

type JoinType string

const (
	InnerJoin JoinType = "JOIN"
	OuterJoin JoinType = "OUTER JOIN"
	LeftJoin  JoinType = "LEFT JOIN"
	RightJoin JoinType = "RIGHT JOIN"
)

type Join struct {
	Schema      string
	OtherSchema string
	Type        JoinType
	JoinClause  Clause
}

func (j Join) Sql(info QueryInfo) (string, []interface{}) {
	clauseSql, vals := j.JoinClause.Sql(info)
	return fmt.Sprintf("%s %s ON %s", j.Type, info.GetJoinSchemas()[j.OtherSchema].Table, clauseSql), vals
}

func (j Join) IsValid(info QueryInfo) bool {
	clauseValid := j.JoinClause.IsValid(info)
	_, otherSchemaValid := info.GetJoinSchemas()[j.OtherSchema]
	return info.GetMainSchema().Name == j.Schema && otherSchemaValid && clauseValid
}
