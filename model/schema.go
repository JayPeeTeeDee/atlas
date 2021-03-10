package model

import (
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"

	"github.com/JayPeeTeeDee/atlas/utils"
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
	Name           string
	ModelType      reflect.Type
	Table          string
	DBNames        []string
	PrimaryFields  []*Field
	Fields         []*Field
	FieldsByName   map[string]*Field
	FieldsByDBName map[string]*Field

	AllFieldNames      *utils.Set
	PrimaryFieldNames  *utils.Set // Used for convenience
	LocationFieldNames *utils.Set
	RegionFieldNames   *utils.Set
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
	HasDefaultValue   bool
	DefaultValue      interface{}
	Schema            Schema
}

func (schema Schema) SetDefaultValues(target interface{}) error {
	vals, err := ParseSingleObject(target, schema)
	if err != nil {
		return err
	}
	for _, field := range schema.Fields {
		if field.HasDefaultValue {
			field.DefaultValue = vals[field.Name]
		}
	}
	return nil
}

func (schema *Schema) ParseField(fieldStruct reflect.StructField) *Field {
	field := &Field{
		Name:              fieldStruct.Name,
		FieldType:         fieldStruct.Type,
		IndirectFieldType: fieldStruct.Type,
		StructField:       fieldStruct,
		Tag:               fieldStruct.Tag,
		TagSettings:       parseTagSetting(fieldStruct.Tag.Get("atlas"), ";"),
	}
	for field.IndirectFieldType.Kind() == reflect.Ptr {
		field.IndirectFieldType = field.IndirectFieldType.Elem()
	}

	fieldValue := reflect.New(field.IndirectFieldType)

	if dbName, ok := field.TagSettings["COLUMN"]; ok {
		field.DBName = dbName
	}

	if val, ok := field.TagSettings["PRIMARYKEY"]; ok && checkTruth(val) {
		field.PrimaryKey = true
	}

	if val, ok := field.TagSettings["AUTOINCREMENT"]; ok && checkTruth(val) {
		field.AutoIncrement = true
	}

	if val, ok := field.TagSettings["NOT NULL"]; ok && checkTruth(val) {
		field.NotNull = true
	}

	if val, ok := field.TagSettings["UNIQUE"]; ok && checkTruth(val) {
		field.Unique = true
	}

	if val, ok := field.TagSettings["DEFAULT"]; ok && checkTruth(val) {
		field.HasDefaultValue = true
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
		Name:               modelType.Name(),
		ModelType:          modelType,
		Table:              toSnakeCase(modelType.Name()),
		FieldsByName:       map[string]*Field{},
		FieldsByDBName:     map[string]*Field{},
		PrimaryFields:      make([]*Field, 0),
		AllFieldNames:      utils.NewSet(),
		PrimaryFieldNames:  utils.NewSet(),
		LocationFieldNames: utils.NewSet(),
		RegionFieldNames:   utils.NewSet(),
	}

	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			schema.Fields = append(schema.Fields, schema.ParseField(fieldStruct))
		}
	}

	for _, field := range schema.Fields {
		if field.DBName == "" && field.DataType != "" {
			field.DBName = toSnakeCase(field.Name)
		}

		if field.DBName != "" {
			if _, ok := schema.FieldsByDBName[field.DBName]; !ok {
				schema.DBNames = append(schema.DBNames, field.DBName)
				schema.FieldsByDBName[field.DBName] = field
				schema.FieldsByName[field.Name] = field
				if field.PrimaryKey {
					schema.PrimaryFields = append(schema.PrimaryFields, field)
					schema.PrimaryFieldNames.Add(field.Name)
				}
				if field.DataType == LocationType {
					schema.LocationFieldNames.Add(field.Name)
				} else if field.DataType == RegionType {
					schema.RegionFieldNames.Add(field.Name)
				}
				schema.AllFieldNames.Add(field.Name)
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

func ParseSingleObject(target interface{}, schema Schema) (map[string]interface{}, error) {
	targetItem := reflect.ValueOf(target)
	if targetItem.Kind() == reflect.Slice {
		return nil, errors.New("Only 1 struct expected")
	} else if targetItem.Kind() == reflect.Struct {
		return parseStruct(targetItem, schema), nil
	} else {
		return nil, errors.New("Unknown input")
	}
}

func parseStruct(targetValue reflect.Value, schema Schema) map[string]interface{} {
	itemRes := make(map[string]interface{})
	for _, field := range schema.Fields {
		itemRes[field.Name] = targetValue.FieldByName(field.Name).Interface()
	}
	return itemRes
}
