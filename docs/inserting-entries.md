# Inserting Entries
## Example
```go
cars := []Car{
    {Id: 1, Location: model.NewLocation(...), OperationZone: ...},
    {Id: 2, Location: model.NewLocation(...), OperationZone: ...},
}
otherCar := Car{Id: 3, Location: ..., OperationZone: ...}

// Can insert single or slice of objects
res, err := db.Create(cars) 
res, err := db.Create(otherCar)

// Can also insert with query API
res, err := db.Model("Car").Create(otherCar)
```

## Specifying All Fields
If all fields of the model are specified, the object can be inserted simply with the following API directly:
```go 
func db.Create(<object or slice of objects>) (sql.Result, error)
```
## Omitting Some Fields
To insert without specify some fields (e.g. auto-incremented fields or fields with default value), 
the query API can be used with the `Omit` method.
```go
car := Car{...}
res, err := db.Model("Car").Omit("Id").Create(car) // Omits Id field in insert query
```

