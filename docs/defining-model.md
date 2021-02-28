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

## Spatial Types
This package provides 2 spatial representations in the `model` subpackage: `Location` and `Region`.
### Location
`Location` is used to denote a point on the WGS84 coordinate system (longitude and latitude).\
It can be instantiated using `model.NewLocation(longitude, latitude)` 
### Region
`Region` is used to denote an area on the WGS84 coordinate system as a list of (longitude, latitude) points.\
It can be instantiated using `model.NewRegion(coords)`.\
_Upcoming feature: Convenience functions for creating regions of basic shapes_
