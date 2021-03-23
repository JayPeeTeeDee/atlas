package query

import (
	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
)

type QueryInfo interface {
	HasSchema(schema string) bool
	HasField(field string) bool
	HasFieldOfType(field string, datatype model.DataType) bool
	GetField(field string) *model.Field
	GetMainSchema() model.Schema
	GetJoinSchemas() map[string]model.Schema
	GetAdapterInfo() adapter.AdapterInfo
}
