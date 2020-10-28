package model

import (
	"reflect"
	"time"
)

type Point struct {
	x float64
	y float64
}

type Location struct {
	point Point
}

type Region struct {
	points []Point
}

type Timestamp struct {
	time *time.Time
}

func NewTimestamp(timestamp time.Time) Timestamp {
	return Timestamp{time: &timestamp}
}

func IsLocation(value reflect.Value) bool {
	if _, ok := value.Interface().(*Location); ok {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(time.Time{})) {
		return true
	} else if value.Type().ConvertibleTo(reflect.TypeOf(&time.Time{})) {
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
