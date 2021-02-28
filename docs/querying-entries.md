# Querying Entries
## Example
```go
// Queries will populate the following variables
requestedCars := make([]Car, 0)
requestedCar := Car{}

// Query for car with id = 1
err := db.Model("Car").Where(query.Equals{Column: "Id", Value: 1}).First(&requestedCar)

// Query only for the location field for car with id = 1 (rest of fields will have zero value)
err := db.Model("Car").Select("Location").Where(query.Equals{Column: "Id", Value: 1}).First(&requestedCar)

// Query for all cars in given region
err := db.Model("Car").CoveredBy(model.NewRegion(...)).All(&requestedCars)
```

## Query Subpackage
Queries (and complex inserts) are formed with the query builder provided by the `query` subpackage. 
The query builder for a model is created using the following method:
```go
db.Model("<model name")
```
The query can then be built upon with the use of _chaining_ methods, before being executed with _terminal_ methods.

### Chaining Methods
| Method Call | SELECT query usage | INSERT query usage | Example |
| --- | --- | --- | --- |
| `Select(...columns)` | Include only these fields in the result | Only specify these fields |`Select("Id", "Location")` |
| `Omit(...columns)` | Exclude these fields in the result | Do not specify these fields | `Omit("OperationZone")` |
| `Limit(count)` | Only return first `count` number of objects | - | `Limit(10)`|
| `Offset(count)` | Return entries after the given offset | - | `Offset(10)` |
| `Where(clause)` | Used to select rows ([See Clauses](#filter-clauses)) | Used to select rows ([See Clauses](#filter-clauses)) | `Where(<clause>)` |

The following chaining methods are provided as convenience 
for spatial [filter clauses](#filter-clauses) in cases where the model has only 1 spatial field (non-ambiguous).

| Method Call | Usage | Example |
| --- | --- | --- |
| `CoveredBy(target)` | Get entries within the target region |`CoveredBy(model.NewRegion(...))` |
| `Covers(target)` | Get entries with region that contain the target spatial object | `Covers(model.NewRegion(...)` |
| `WithinRangeOf(targets, range)` | Get entries that are within `range` meters from _any_ of the targets | `WithinRangeOf([]model.Location{...}, 10)`|
| `HasWithinRange(targets, range)` | Get entries that have _all_ targets within `range` meters | - | `HasWithinRange([]model.Location{...}, 10)` |

### Terminal Methods
Results of the query (entries or count) are populated into the pointers passed into the functions. 
For INSERT queries, terminal methods will also return any error encountered when building or executing the query.

| Method Signature | SELECT query usage | Example |
| --- | --- | --- |
| `Count(count *int) error` | Counts the number of entries | `Count(&count)` |
| `First(response interface{}) error` | Get only the first entry | `First(&car)` |
| `All(response interface{}) error` | Get all entries | `All(&cars)`|

## Filter Clauses
Filter clauses are used to specify conditions for the `Where` chaining method. 
These are structs that can be instantiated and passed as the paramter to the `Where` method.

### Basic Conditional Clauses
The following clauses are used to compare entries against values
#### Equal
Get entries where model.Column is equal to Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types, `Location` and `Region`)
- Example:
  - `Equal{Column: "Id", Value: 1}`
  - `Equal{Column: "Location", Value: model.NewLocation(...)}`
  
#### NotEqual
Get entries where model.Column is not equal to Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types, `Location` and `Region`)
- Example:
  - `NotEqual{Column: "Id", Value: 1}`
  - `NotEqual{Column: "Location", Value: model.NewLocation(...)}`
  
#### GreaterThan
Get entries where model.Column is greater than Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types)
- Example:
  - `GreaterThan{Column: "Id", Value: 1}`
  
#### GreaterThanOrEqual
Get entries where model.Column is greater than or equal to Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types)
- Example:
   - `GreaterThanOrEqual{Column: "Id", Value: 1}`

#### LessThan
Get entries where model.Column is less than Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types)
- Example:
  - `LessThan{Column: "Id", Value: 1}`
  
#### LessThanOrEqual
Get entries where model.Column is less than or equal to Value
- Parameters:
  - Column: Field name of model to compare
  - Value: Value to check against (supports golang primitive types)
- Example:
  - `LessThanOrEqual{Column: "Id", Value: 1}`

#### Like
Get entries where model.Column matches the given pattern (only for strings)
- Parameters:
  - Column: Field name of model to compare (of string type)
  - Value: Pattern to match against
- Example:
  - `Like{Column: "Name", Value: "jo_"}`
  
#### NotLike
Get entries where model.Column does not match the given pattern (only for strings)
- Parameters:
  - Column: Field name of model to compare (of string type)
  - Value: Pattern to match against
- Example:
  - `NotLike{Column: "Name", Value: "jo_"}`
 
### Spatial Conditional Clauses
The following clauses are used to compare spatial fields against spatial values

#### CoveredBy
Get entries where model.Column is covered by the target spatial object (`Location` or `Region`)
- Parameters:
  - Column: Field name of model to compare
  - Target: Spatial object to compare against (`Location` or `Region`)
- Example:
  - `CoveredBy{Column: "OperationZone", Target: model.NewRegion(...)}`
 
#### Covers
Get entries where model.Column covers the target spatial object (`Location` or `Region`)
- Parameters:
  - Column: Field name of model to compare
  - Target: Spatial object to compare against (`Location` or `Region`)
- Example:
  - `Covers{Column: "OperationZone", Target: model.NewRegion(...)}`
  
#### WithinRangeOf
Get entries where model.Column is within range of _any_ of the target objects (`Location` or `Region`)
- Parameters:
  - Column: Field name of model to compare
  - Targets: Spatial objects to compare against (`Location` or `Region`)
  - Range: Distance in meters
- Example:
  - `WithinRangeOf{Column: "OperationZone", Targets: []model.Location{...}}`
 
#### HasWithinRange
Get entries where model.Column is within range of _all_ the target objects (`Location` or `Region`)
- Parameters:
  - Column: Field name of model to compare
  - Targets: Spatial objects to compare against (`Location` or `Region`)
  - Range: Distance in meters
- Example:
  - `HasWithinRange{Column: "OperationZone", Targets: []model.Location{...}}`
  
### Combination Clauses
The following clauses are used to combine different conditional clauses together.

#### And
Get entries that satisfy all clauses
- Example:
  - `And{Equal{...}, Covers{...}}`
 
#### Or
Get entries that satisfy any of the clauses
- Example:
  - `Or{Equal{...}, Covers{...}}`
