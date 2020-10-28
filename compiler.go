package atlas

import (
	"github.com/JayPeeTeeDee/atlas/query"
	"github.com/Masterminds/squirrel"
)

func CompileSQL(builder query.Builder) (string, []interface{}) {
	if builder.QueryType == query.SelectQuery {
		selection := "*"
		if builder.IsCount {
			selection = "COUNT(*)"
		}
		// TODO: Adjust placeholder format by db adapter
		statementBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select(selection).From(builder.TableName)

		for _, condition := range builder.Conditions {
			conditionString := condition.Column + " " + condition.Expr.GetCondition() + " ?"
			statementBuilder = statementBuilder.Where(conditionString, condition.Expr.GetValue())
		}

		if builder.Limit > 0 {
			statementBuilder = statementBuilder.Limit(builder.Limit)
		}

		if builder.Offset > 0 {
			statementBuilder = statementBuilder.Offset(builder.Offset)
		}

		sql, args, _ := statementBuilder.ToSql()
		return sql, args
	}
	return "", nil
}
