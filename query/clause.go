package query

import (
	"fmt"
	"strings"
)

type Clause interface {
	Condition() string
	Sql() (string, []interface{})
}

type GreaterThan struct {
	Column string
	Value  string
}

func (e GreaterThan) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s > ?", e.Column), []interface{}{e.Value}
}

func (e GreaterThan) Condition() string {
	return ">"
}

type LessThan struct {
	Column string
	Value  string
}

func (e LessThan) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s < ?", e.Column), []interface{}{e.Value}
}

func (e LessThan) Condition() string {
	return "<"
}

type Equal struct {
	Column string
	Value  string
}

func (e Equal) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s = ?", e.Column), []interface{}{e.Value}
}

func (e Equal) Condition() string {
	return "="
}

type GreaterThanOrEqual struct {
	Column string
	Value  string
}

func (e GreaterThanOrEqual) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s >= ?", e.Column), []interface{}{e.Value}
}

func (e GreaterThanOrEqual) Condition() string {
	return ">="
}

type LessThanOrEqual struct {
	Column string
	Value  string
}

func (e LessThanOrEqual) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s <= ?", e.Column), []interface{}{e.Value}
}

func (e LessThanOrEqual) Condition() string {
	return "<="
}

type NotEqual struct {
	Column string
	Value  string
}

func (e NotEqual) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s <> ?", e.Column), []interface{}{e.Value}
}

func (e NotEqual) Condition() string {
	return "<>"
}

type Like struct {
	Column string
	Value  string
}

func (e Like) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s LIKE ?", e.Column), []interface{}{e.Value}
}

func (e Like) Condition() string {
	return "LIKE"
}

type NotLike struct {
	Column string
	Value  string
}

func (e NotLike) Sql() (string, []interface{}) {
	return fmt.Sprintf("%s NOT LIKE ?", e.Column), []interface{}{e.Value}
}

func (e NotLike) Condition() string {
	return "NOT LIKE"
}

type Or []Clause

func (e Or) Sql() (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql()
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

func (e And) Sql() (string, []interface{}) {
	sql := strings.Builder{}
	values := make([]interface{}, 0)
	for i, clause := range e {
		clauseSql, clauseVals := clause.Sql()
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
