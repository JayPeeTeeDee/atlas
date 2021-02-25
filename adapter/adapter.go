package adapter

import "database/sql"

type PlaceholderStyle string

const QuestionPlaceholder PlaceholderStyle = "QUESTION_PLACEHOLDER"
const DollarPlaceholder PlaceholderStyle = "DOLLAR_PLACEHOLDER"

type SpatialExtension string

const PostGisExtension SpatialExtension = "POSTGIS_EXTENSION"

type Adapter interface {
	Connect(dsn string) error
	Disconnect() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Placeholder() PlaceholderStyle
	SpatialType() SpatialExtension
}
