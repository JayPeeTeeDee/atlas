package query

import (
	"fmt"

	"github.com/JayPeeTeeDee/atlas/adapter"
	"github.com/JayPeeTeeDee/atlas/model"
)

type Order interface {
	IsDescending() bool
	IsValid(fields map[string]*model.Field) bool
	Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{})
}

type ColumnOrder struct {
	Column     string
	Descending bool
}

func (c ColumnOrder) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	sql := fmt.Sprintf("%s ", fields[c.Column].DBName)
	if c.Descending {
		sql += "DESC"
	} else {
		sql += "ASC"
	}
	return sql, []interface{}{}
}

func (c ColumnOrder) IsValid(fields map[string]*model.Field) bool {
	col, ok := fields[c.Column]
	if !ok {
		return ok
	} else {
		return col.DataType != model.LocationType && col.DataType != model.RegionType
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

func (s SpatialOrder) Sql(fields map[string]*model.Field, spatialType adapter.SpatialExtension) (string, []interface{}) {
	sql := ""
	switch spatialType {
	case adapter.PostGisExtension:
		sql += fmt.Sprintf("%s::geometry <#> ST_GeomFromGeoJSON(?) ", fields[s.Column].DBName)
	default:
		sql += fmt.Sprintf("%s::geometry <#> ST_GeomFromGeoJSON(?) ", fields[s.Column].DBName)
	}
	if s.Descending {
		sql += "DESC"
	} else {
		sql += "ASC"
	}
	return sql, []interface{}{s.Target}
}

func (s SpatialOrder) IsValid(fields map[string]*model.Field) bool {
	col, ok := fields[s.Column]
	if !ok || (col.DataType != model.LocationType && col.DataType != model.RegionType) {
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
