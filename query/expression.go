package query

type Expression interface {
	GetCondition() string
	GetValue() string
}

type GreaterThan struct {
	Value string
}

func (e GreaterThan) GetValue() string {
	return e.Value
}

func (e GreaterThan) GetCondition() string {
	return ">"
}

type LessThan struct {
	Value string
}

func (e LessThan) GetValue() string {
	return e.Value
}

func (e LessThan) GetCondition() string {
	return "<"
}

type Equal struct {
	Value string
}

func (e Equal) GetValue() string {
	return e.Value
}

func (e Equal) GetCondition() string {
	return "="
}

type GreaterThanOrEqual struct {
	Value string
}

func (e GreaterThanOrEqual) GetValue() string {
	return e.Value
}

func (e GreaterThanOrEqual) GetCondition() string {
	return ">="
}

type LessThanOrEqual struct {
	Value string
}

func (e LessThanOrEqual) GetValue() string {
	return e.Value
}

func (e LessThanOrEqual) GetCondition() string {
	return "<="
}

type NotEqual struct {
	Value string
}

func (e NotEqual) GetValue() string {
	return e.Value
}

func (e NotEqual) GetCondition() string {
	return "<>"
}
