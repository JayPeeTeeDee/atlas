package query

type Aggregation interface {
	GetAggregation() string
}

type Selection struct {
	column      string
	aggregation Aggregation
}

type Order struct {
	Column     string
	Descending bool
}

type Type string

const (
	SelectQuery Type = "SelectQueryType"
	InsertQuery Type = "InsertQueryType"
)

type Builder struct {
	TableName  string
	selections []Selection
	Clauses    []Clause
	Orders     []Order
	Limit      uint64
	Offset     uint64
	QueryType  Type
	IsCount    bool
	// TODO: join, having, groupby, returning
}

func (b *Builder) Where(clause Clause) *Builder {
	b.Clauses = append(b.Clauses, clause)
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
