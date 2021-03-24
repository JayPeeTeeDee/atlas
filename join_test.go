package atlas

import (
	"testing"

	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
)

type CarTest struct {
	CarId          int `atlas:"primarykey;autoincrement"`
	Location       model.Location
	DesignatedZone string
}

type ZoneTest struct {
	ZoneId int `atlas:"primarykey;autoincrement"`
	Name   string
	Region model.Region
}

func setup() (*Database, error) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		return nil, err
	}

	err = db.RegisterModel(CarTest{})
	if err != nil {
		tearDown(db)
		return nil, err
	}
	err = db.RegisterModel(ZoneTest{})
	if err != nil {
		tearDown(db)
		return nil, err
	}

	err = db.CreateTable("CarTest", true)
	if err != nil {
		tearDown(db)
		return nil, err
	}
	err = db.CreateTable("ZoneTest", true)
	if err != nil {
		tearDown(db)
		return nil, err
	}
	return db, nil
}

func tearDown(db *Database) {
	db.Execute("DROP TABLE IF EXISTS car_test;")
	db.Execute("DROP TABLE IF EXISTS zone_test;")
	db.Disconnect()
}

func populateRows(db *Database) error {
	zones := []ZoneTest{
		{Name: "North", Region: model.NewRectRegion(-10, 10, 10, 20)},
		{Name: "South", Region: model.NewRectRegion(-10, 10, -20, -10)},
		{Name: "East", Region: model.NewRectRegion(10, 20, -10, 10)},
		{Name: "West", Region: model.NewRectRegion(-20, -10, -10, 10)},
	}

	_, err := db.Model("ZoneTest").Omit("ZoneId").Create(zones)
	if err != nil {
		return err
	}

	cars := []CarTest{
		// In correct zone
		{Location: model.NewLocation(5, 11), DesignatedZone: "North"},
		{Location: model.NewLocation(5, -11), DesignatedZone: "South"},
		{Location: model.NewLocation(15, 5), DesignatedZone: "East"},
		{Location: model.NewLocation(-20, 7), DesignatedZone: "West"},
		// In wrong zone
		{Location: model.NewLocation(5, 5), DesignatedZone: "North"},
		{Location: model.NewLocation(5, 20), DesignatedZone: "South"},
		{Location: model.NewLocation(9, 5), DesignatedZone: "East"},
		{Location: model.NewLocation(-30, 7), DesignatedZone: "West"},
	}

	_, err = db.Model("CarTest").Omit("CarId").Create(cars)
	if err != nil {
		return err
	}
	return nil
}

func TestSimpleJoin(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("Failed to set up: %s\n", err.Error())
	}

	err = populateRows(db)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := []CarTest{}
	err = db.Model("CarTest").Join("ZoneTest", query.Equal{Column: "CarTest.DesignatedZone", OtherColumn: "ZoneTest.Name"}).Select("CarId", "CarTest.Location").Where(query.Equal{Column: "ZoneTest.Name", Value: "North"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("Query is wrong")
	} else if !res[0].Location.IsEqual(model.NewLocation(5, 11)) && !res[1].Location.IsEqual(model.NewLocation(5, 11)) {
		t.Errorf("Query is wrong")
	}

	tearDown(db)
}

func TestSpatialJoin(t *testing.T) {
	db, err := setup()
	if err != nil {
		t.Errorf("Failed to set up: %s\n", err.Error())
	}

	err = populateRows(db)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := []CarTest{}
	fullQuery := db.Model("CarTest").Join("ZoneTest", query.CoveredBy{Column: "CarTest.Location", TargetColumn: "ZoneTest.Region"}).Select("CarId", "DesignatedZone").Where(query.Equal{Column: "ZoneTest.Name", Value: "North"})
	err = fullQuery.All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}
	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	err = fullQuery.Where(query.NotEqual{Column: "ZoneTest.Name", OtherColumn: "CarTest.DesignatedZone"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}
	if len(res) != 1 || res[0].DesignatedZone != "South" {
		t.Errorf("Query is wrong")
	}

	tearDown(db)
}
