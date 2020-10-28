package query

type Aggregation interface {
	GetAggregation() string
}

// type Selection struct {
// 	column      string
// 	aggregation Aggregation
// }

type Order struct {
	Column     string
	Descending bool
}

type Condition struct {
	Column string
	Expr   Expression
}

type Type string

const (
	SelectQuery Type = "SelectQueryType"
	CountQuery  Type = "SelectQueryType"
	ExecQuery   Type = "SelectQueryType"
)

type Builder struct {
	TableName string
	// selections []Selection
	Conditions []Condition
	Orders     []Order
	Limit      uint64
	Offset     uint64
	QueryType  Type
	IsCount    bool
	// TODO: join, having, groupby, returning
}

func (b *Builder) Where(column string, expression Expression) *Builder {
	condition := Condition{Column: column, Expr: expression}
	b.Conditions = append(b.Conditions, condition)
	return b
}

func (b *Builder) OrderBy(column string, desc bool) *Builder {
	// TODO: check that table has column
	order := Order{
		Column:     column,
		Descending: desc,
	}

	b.Orders = append(b.Orders, order)
	return b
}
