package model

import (
	"testing"
)

type TestStruct struct {
	Name   string `atlas:"primarykey"`
	Region Region
}

func TestSchema(t *testing.T) {
	schema, err := Parse(TestStruct{})
	if err != nil {
		t.Error(err)
	}
	if schema.Fields[1].DataType != RegionType {
		t.Errorf("Not a region")
	}
	t.Log(schema.Fields[1].DataType)
	t.Log(schema.Table)
}
