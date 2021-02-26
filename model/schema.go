package model

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
)

type DataType string

var ErrUnsupportedDataType = errors.New("unsupported data type")

const (
	Bool          DataType = "bool"
	Int           DataType = "int"
	Uint          DataType = "uint"
	Float         DataType = "float"
	String        DataType = "string"
	Time          DataType = "time"
	Bytes         DataType = "bytes"
	LocationType  DataType = "location"
	RegionType    DataType = "region"
	TimestampType DataType = "timestamp"
)

type Schema struct {
	Name                string
	ModelType           reflect.Type
	Table               string
	DBNames             []string
	PrimaryFields       []*Field
	PrimaryFieldDBNames []string
	Fields              []*Field
	FieldsByName        map[string]*Field
	FieldsByDBName      map[string]*Field
}

type Field struct {
	Name              string
	DBName            string
	FieldType         reflect.Type
	IndirectFieldType reflect.Type
	DataType          DataType
	StructField       reflect.StructField
	Tag               reflect.StructTag
	TagSettings       map[string]string
	PrimaryKey        bool
	AutoIncrement     bool
	NotNull           bool
	Unique            bool
	Schema            Schema
}

func (schema *Schema) ParseField(fieldStruct reflect.StructField) *Field {
	field := &Field{
		Name:              fieldStruct.Name,
		FieldType:         fieldStruct.Type,
		IndirectFieldType: fieldStruct.Type,
		StructField:       fieldStruct,
		Tag:               fieldStruct.Tag,
		TagSettings:       ParseTagSetting(fieldStruct.Tag.Get("atlas"), ";"),
	}
	for field.IndirectFieldType.Kind() == reflect.Ptr {
		field.IndirectFieldType = field.IndirectFieldType.Elem()
	}

	fieldValue := reflect.New(field.IndirectFieldType)

	if dbName, ok := field.TagSettings["COLUMN"]; ok {
		field.DBName = dbName
	}

	if val, ok := field.TagSettings["PRIMARYKEY"]; ok && CheckTruth(val) {
		field.PrimaryKey = true
	}

	if val, ok := field.TagSettings["AUTOINCREMENT"]; ok && CheckTruth(val) {
		field.AutoIncrement = true
	}

	if val, ok := field.TagSettings["NOT NULL"]; ok && CheckTruth(val) {
		field.NotNull = true
	}

	if val, ok := field.TagSettings["UNIQUE"]; ok && CheckTruth(val) {
		field.Unique = true
	}

	switch reflect.Indirect(fieldValue).Kind() {
	case reflect.Bool:
		field.DataType = Bool
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		field.DataType = Int
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		field.DataType = Uint
	case reflect.Float32, reflect.Float64:
		field.DataType = Float
	case reflect.String:
		field.DataType = String
	case reflect.Struct:
		if IsLocation(fieldValue) {
			field.DataType = LocationType
		} else if IsRegion(fieldValue) {
			field.DataType = RegionType
		} else if IsTimestamp(fieldValue) {
			field.DataType = TimestampType
		}
	case reflect.Array, reflect.Slice:
		if reflect.Indirect(fieldValue).Type().Elem() == reflect.TypeOf(uint8(0)) {
			field.DataType = Bytes
		}
	}

	if val, ok := field.TagSettings["TYPE"]; ok {
		switch DataType(strings.ToLower(val)) {
		case Bool, Int, Uint, Float, String, Time, Bytes:
			field.DataType = DataType(strings.ToLower(val))
		default:
			field.DataType = DataType(val)
		}
	}

	return field
}

func Parse(target interface{}) (*Schema, error) {
	if target == nil {
		return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, target)
	}

	modelType := reflect.ValueOf(target).Type()
	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			return nil, fmt.Errorf("%w: %+v", ErrUnsupportedDataType, target)
		}
		return nil, fmt.Errorf("%w: %v.%v", ErrUnsupportedDataType, modelType.PkgPath(), modelType.Name())
	}

	schema := &Schema{
		Name:           modelType.Name(),
		ModelType:      modelType,
		Table:          ToSnakeCase(modelType.Name()),
		FieldsByName:   map[string]*Field{},
		FieldsByDBName: map[string]*Field{},
	}

	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			schema.Fields = append(schema.Fields, schema.ParseField(fieldStruct))
		}
	}

	for _, field := range schema.Fields {
		if field.DBName == "" && field.DataType != "" {
			field.DBName = ToSnakeCase(field.Name)
		}

		if field.DBName != "" {
			// nonexistence or shortest path or first appear prioritized if has permission
			if v, ok := schema.FieldsByDBName[field.DBName]; !ok {
				if _, ok := schema.FieldsByDBName[field.DBName]; !ok {
					schema.DBNames = append(schema.DBNames, field.DBName)
				}
				schema.FieldsByDBName[field.DBName] = field
				schema.FieldsByName[field.Name] = field

				if v != nil && v.PrimaryKey {
					for idx, f := range schema.PrimaryFields {
						if f == v {
							schema.PrimaryFields = append(schema.PrimaryFields[0:idx], schema.PrimaryFields[idx+1:]...)
						}
					}
				}

				if field.PrimaryKey {
					schema.PrimaryFields = append(schema.PrimaryFields, field)
				}
			}
		}

		if _, ok := schema.FieldsByName[field.Name]; !ok {
			schema.FieldsByName[field.Name] = field
		}
	}

	return schema, nil
}

func ParseType(target interface{}) (string, error) {
	if target == nil {
		return "", fmt.Errorf("%w: %+v", ErrUnsupportedDataType, target)
	}
	modelType := reflect.ValueOf(target).Type()
	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	return modelType.Name(), nil
}

func ParseObject(target interface{}, schema Schema) ([]map[string]interface{}, error) {
	targetItem := reflect.ValueOf(target)
	res := make([]map[string]interface{}, 0)
	if targetItem.Kind() == reflect.Slice {
		for i := 0; i < targetItem.Len(); i++ {
			item := targetItem.Index(i)
			res = append(res, parseStruct(item, schema))
		}
	} else if targetItem.Kind() == reflect.Struct {
		res = append(res, parseStruct(targetItem, schema))
	}

	return res, nil
}

func parseStruct(targetValue reflect.Value, schema Schema) map[string]interface{} {
	itemRes := make(map[string]interface{})
	for _, field := range schema.Fields {
		itemRes[field.Name] = targetValue.FieldByName(field.Name).Interface()
	}
	return itemRes
}
