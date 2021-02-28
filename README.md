Atlas -- Optimised ORM for building location-based services with spatio-temporal data

# Features
- Defining models with region and location
- Querying spatial models with PostGIS plugin

# Example Usage
```go
package main

import (
	"fmt"
	"os"

	"github.com/JayPeeTeeDee/atlas"
	"github.com/JayPeeTeeDee/atlas/model"
)

type Car struct {
	Id       int `atlas:"primarykey"`
	Location model.Location
	Brand    string
	Model    string
}

func main() {
	db, err := atlas.ConnectWithDSN(atlas.DBType_Postgres, "postgresql://username:password@localhost/database")
	if err != nil {
		fmt.Fprint(os.Stderr, "Unable to connect to database")
		os.Exit(1)
	}
	defer db.Disconnect()

	// Upcoming feature: Table creation
	_, err = db.Execute("CREATE TABLE IF NOT EXISTS car(id int PRIMARY KEY, location geography, brand varchar(100), model varchar(100));")
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	db.RegisterModel(&Car{})

	cars := []Car{
		{Id: 1, Location: model.NewLocation(103.81, 1.30), Brand: "Toyota", Model: "Corolla Altis"},
		{Id: 2, Location: model.NewLocation(101.97, 4.3), Brand: "Mitsubishi", Model: "Lancer"},
	}
	_, err = db.Create(cars)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to insert cars")
		os.Exit(1)
	}

	carsInRegion := make([]Car, 0)
	err = db.Model("Car").CoveredBy(model.NewRegion(
		[][]float64{
			{100, 10},
			{102, 10},
			{102, 0},
			{100, 0},
			{100, 10},
		},
	)).All(&carsInRegion)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to query cars")
		os.Exit(1)
	}

	fmt.Printf("There are %d cars in the region\n", len(carsInRegion))
	fmt.Printf("The first car is at %v\n", carsInRegion[0].Location)
}

```

# Documentation
- [Defining Model](https://github.com/JayPeeTeeDee/atlas/blob/master/docs/defining-model.md)
- [Inserting Entries](https://github.com/JayPeeTeeDee/atlas/blob/master/docs/inserting-entries.md)
- [Querying Entries](https://github.com/JayPeeTeeDee/atlas/blob/master/docs/querying-entries.md)


