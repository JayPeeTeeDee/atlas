package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	geojson "github.com/paulmach/go.geojson"
)

type SpatialObject interface {
	IsLocation() bool
	IsRegion() bool
}

type Location struct {
	point *geojson.Geometry
}

func NewLocation(lon float64, lat float64) Location {
	return Location{geojson.NewPointGeometry([]float64{lon, lat})}
}

func (l *Location) IsEqual(other Location) bool {
	first_val, first_err := l.point.Value()
	other_val, other_err := other.point.Value()
	return first_err == nil && other_err == nil && reflect.DeepEqual(first_val, other_val)
}

func (l *Location) Scan(value interface{}) error {
	l.point = &geojson.Geometry{}
	if value == nil {
		return nil
	}
	err := l.point.Scan(value)
	if err != nil {
		return err
	}
	if !l.point.IsPoint() {
		return errors.New("Invalid location type from database")
	}
	return nil
}

func (l Location) Value() (driver.Value, error) {
	if !l.point.IsPoint() {
		return nil, errors.New("Invalid location representation")
	}
	return l.point.Value()
}

func (l Location) String() string {
	return fmt.Sprintf("(%f, %f)", l.point.Point[0], l.point.Point[1])
}

func (l Location) IsLocation() bool {
	return true
}

func (l Location) IsRegion() bool {
	return false
}

type Region struct {
	polygon *geojson.Geometry
}

func NewRegion(coords [][]float64) Region {
	poly := make([][][]float64, 0)
	poly = append(poly, coords)
	return Region{geojson.NewPolygonGeometry(poly)}
}

func (r *Region) IsEqual(other Region) bool {
	first_val, first_err := r.polygon.Value()
	other_val, other_err := other.polygon.Value()
	return first_err == nil && other_err == nil && reflect.DeepEqual(first_val, other_val)
}

func (r *Region) Scan(value interface{}) error {
	r.polygon = &geojson.Geometry{}
	if value == nil {
		return nil
	}
	err := r.polygon.Scan(value)
	if err != nil {
		return err
	}
	if !r.polygon.IsPolygon() {
		return errors.New("Invalid region type from database")
	}
	return nil
}

func (r Region) Value() (driver.Value, error) {
	if !r.polygon.IsPolygon() {
		return nil, errors.New("Invalid region representation")
	}
	return r.polygon.Value()
}

func (r Region) String() string {
	if len(r.polygon.Polygon) == 0 {
		return ""
	}
	strs := make([]string, len(r.polygon.Polygon[0]))
	for i, point := range r.polygon.Polygon[0] {
		strs[i] = fmt.Sprintf("(%f, %f)", point[0], point[1])
	}
	return fmt.Sprintf("(%s)", strings.Join(strs, ","))
}

func (r Region) IsLocation() bool {
	return false
}

func (r Region) IsRegion() bool {
	return true
}

type Timestamp struct {
	time *time.Time
}

func (t *Timestamp) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		t.time = &v
		return nil
	default:
		return errors.New("Invalid timestamp type from database")
	}
}

func (t Timestamp) Value() (driver.Value, error) {
	return *(t.time), nil
}

func NewTimestamp(timestamp time.Time) Timestamp {
	return Timestamp{time: &timestamp}
}

func IsLocation(value reflect.Value) bool {
	if _, ok := value.Interface().(*Location); ok {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(Location{})) {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(&Location{})) {
		return true
	}
	return false
}

func IsRegion(value reflect.Value) bool {
	if _, ok := value.Interface().(*Region); ok {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(Region{})) {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(&Region{})) {
		return true
	}
	return false
}

func IsTimestamp(value reflect.Value) bool {
	if _, ok := value.Interface().(*Timestamp); ok {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(Timestamp{})) {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(&Timestamp{})) {
		return true
	}
	return false
}
