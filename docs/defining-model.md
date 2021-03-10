# Defining Model
## Example
```go
type Car struct {
    Id            int `atlas:"primarykey"`
    Location      model.Location
    OperationZone model.Region
}
```

## Field Tags
Tags are used to indicate properties of specific fields of the model. Tags used in this package are prefixed with `atlas`.\
They follow the format of `atlas:"<tag name>"` for key only tags and `atlas:"<tag name>:<value>"` for key-value tags.\
For multiple tags, they should be separated by semicolons (e.g. `atlas:"<tag1>;<tag2>:<tag2val>"`)

| Tag Name | Type | Description |
| --- | --- | --- |
| `primarykey` | key |Indicates that this field is part of the primary key of the model |
| `column` | key-value | Indicates the column name of the field in the database |
| `not null` | key | Indicates that this field should not be null |
| `unique` | key | Indicates that this field should have unique values |
| `default` | key | Indicates that this field should have default value (defined during registration) |

## Spatial Types
This package provides 2 spatial representations in the `model` subpackage: `Location` and `Region`.
### Location
`Location` is used to denote a point on the WGS84 coordinate system (longitude and latitude).\
It can be instantiated using `model.NewLocation(longitude, latitude)` 
### Region
`Region` is used to denote an area on the WGS84 coordinate system as a list of (longitude, latitude) points.\
It can be instantiated using `model.NewRegion(coords)`.\
If the region is a regular rectangle, the following convenience function can be used: 
`model.NewRectRegion(minLon, maxLon, minLat, maxLat)`.

## Registering Model with Atlas
Once the model struct has been defined, it needs to be registered for Atlas to recognise it for queries:
```go
func (d *Database) RegisterModel(<object>) error
// e.g. db.RegisterModel(Car{})
```
Any default values for fields (specified using the `default` field tag) will need to be specified in the object passed into `RegisterModel`:
```go
type Car struct {
    Id            int `atlas:"primarykey"`
    Location      model.Location `atlas:"default"`
    OperationZone model.Region
}

err := db.RegisterModel(Car{Location: model.NewLocation(103.82, 1.35)})
```

## Creating Table on DBMS with Atlas
After registering the model with Atlas, it can automatically create tables in the database for you:
```go
func (d *Database) CreateTable(schemaName string, ifNotExists bool) error
// e.g. db.CreateTable("Car", true)
```