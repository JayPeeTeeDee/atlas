package query

import (
	"fmt"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
)

type Order interface {
	IsDescending() bool
	IsValid(info QueryInfo) bool
	Sql(info QueryInfo) (string, []interface{})
}

type ColumnOrder struct {
	Column     string
	Descending bool
}

func (c ColumnOrder) Sql(info QueryInfo) (string, []interface{}) {
	sql := fmt.Sprintf("%s ", info.GetField(c.Column).DBName)
	if c.Descending {
		sql += "DESC"
	} else {
		sql += "ASC"
	}
	return sql, []interface{}{}
}

func (c ColumnOrder) IsValid(info QueryInfo) bool {
	field := info.GetField(c.Column)
	if field == nil {
		return false
	} else {
		return field.DataType != model.LocationType && field.DataType != model.RegionType
	}
}

func (c ColumnOrder) IsDescending() bool {
	return c.Descending
}

type SpatialOrder struct {
	Column     string
	Descending bool
	Target     model.SpatialObject
}

func (s SpatialOrder) Sql(info QueryInfo) (string, []interface{}) {
	spatialType := info.GetAdapterInfo().SpatialType()
	sql := ""
	switch spatialType {
	case adapter.PostGisExtension:
		sql += fmt.Sprintf("%s::geometry <#> ST_GeomFromGeoJSON(?) ", info.GetField(s.Column).DBName)
	default:
		sql += fmt.Sprintf("%s::geometry <#> ST_GeomFromGeoJSON(?) ", info.GetField(s.Column).DBName)
	}
	if s.Descending {
		sql += "DESC"
	} else {
		sql += "ASC"
	}
	return sql, []interface{}{s.Target}
}

func (s SpatialOrder) IsValid(info QueryInfo) bool {
	field := info.GetField(s.Column)
	if field == nil || (field.DataType != model.LocationType && field.DataType != model.RegionType) {
		return false
	} else {
		switch s.Target.(type) {
		case model.Location, model.Region:
			return true
		default:
			return false
		}
	}
}

func (s SpatialOrder) IsDescending() bool {
	return s.Descending
}
