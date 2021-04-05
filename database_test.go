package atlas

import (
	"testing"

	"github.com/JayPeeTeeDee/atlas/model"
	"github.com/JayPeeTeeDee/atlas/query"
)

func TestConnection(t *testing.T) {
	_, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua@localhost/johnphua")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}
}

func TestRawQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	_, err = db.Execute("CREATE TABLE test (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO test VALUES (1, 1), (2, 2)")
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	rows, err := db.Query("SELECT * FROM test;")
	if err != nil {
		t.Errorf("Failed to retrieve rows: %s\n", err.Error())
	}
	count := 0
	for rows.Next() {
		count += 1
	}
	rows.Close()
	if count != 2 {
		t.Errorf("Rows do not match up")
	}

	_, err = db.Execute("DROP TABLE test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

type ModelTest struct {
	A int
	B int `atlas:"default"`
}

func TestModelSingleQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO model_test VALUES (1, 1), (2, 2), (1, 2)")
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := ModelTest{}
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).First(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if res.A != 1 {
		t.Errorf("Query is wrong")
	}

	count := 0
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).Count(&count)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}
	if count != 2 {
		t.Errorf("Query is wrong")
	}

	res = ModelTest{}
	err = db.Model("ModelTest").Where(query.And{
		query.Equal{Column: "ModelTest.A", Value: "1"},
		query.Equal{Column: "B", Value: "1"},
	}).First(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if res.A != 1 || res.B != 1 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestModelSelectQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO model_test VALUES (1, 1), (2, 2), (1, 2)")
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := ModelTest{}
	err = db.Model("ModelTest").Select("B").Where(query.Equal{Column: "A", Value: "2"}).First(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if res.A != 0 || res.B != 2 {
		t.Errorf("Query is wrong")
	}

	other := make([]ModelTest, 0)
	err = db.Model("ModelTest").Select("ModelTest.A").Where(query.Equal{Column: "A", Value: "1"}).All(&other)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(other) != 2 {
		t.Errorf("Query is wrong")
	}
	for _, model := range other {
		if model.A != 1 || model.B != 0 {
			t.Errorf("Query is wrong")
		}
	}

	other = make([]ModelTest, 0)
	err = db.Model("ModelTest").Select("A", "B").Where(query.Equal{Column: "A", Value: "1"}).All(&other)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(other) != 2 {
		t.Errorf("Query is wrong")
	}
	for _, model := range other {
		if model.A != 1 || model.B == 0 {
			t.Errorf("Query is wrong")
		}
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestModelAllQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO model_test VALUES (1, 1), (2, 2), (1, 2)")
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	res = make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.And{
		query.Equal{Column: "A", Value: "1"},
		query.Equal{Column: "B", Value: "1"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestModelSingleInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	val := ModelTest{A: 1, B: 1}

	_, err = db.Create(val)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong")
	}

	res = make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.And{
		query.Equal{Column: "A", Value: "1"},
		query.Equal{Column: "B", Value: "2"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 0 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}
func TestModelMultiInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []ModelTest{
		{A: 1, B: 1},
		{A: 2, B: 2},
		{A: 1, B: 2},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	res = make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.And{
		query.Equal{Column: "A", Value: "1"},
		query.Equal{Column: "B", Value: "1"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestModelSelectInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{B: 1})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []ModelTest{
		{A: 1},
		{A: 2},
		{A: 1},
	}

	_, err = db.Model("ModelTest").Select("A").Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	res = make([]ModelTest, 0)
	err = db.Model("ModelTest").Select("A").Where(query.And{
		query.Equal{Column: "A", Value: "2"},
		query.Equal{Column: "B", Value: "1"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 || res[0].B != 0 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestModelOmitInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(ModelTest{B: 1})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("ModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []ModelTest{
		{A: 1},
		{A: 2},
		{A: 1},
	}

	_, err = db.Model("ModelTest").Omit("B").Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "A", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	res = make([]ModelTest, 0)
	err = db.Model("ModelTest").Omit("B").Where(query.And{
		query.Equal{Column: "A", Value: "2"},
		query.Equal{Column: "B", Value: "1"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 || res[0].B != 0 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

type SpatialModelTest struct {
	A int `atlas:"primarykey"`
	B model.Location
	C model.Region
}

func TestSpatialModelInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(SpatialModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("SpatialModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SpatialModelTest{
		{
			A: 1,
			B: model.NewLocation(-60.0, 0.0),
			C: model.NewRectRegion(-100.0, 100.0, -100.0, 100.0),
		},
		{
			A: 2,
			B: model.NewLocation(30.0, 10.0),
			C: model.NewRegion([][]float64{
				{-100.0, -100.0},
				{-100.0, 100.0},
				{100.0, 100.0},
				{100.0, -100.0},
			}),
		},
		{
			A: 3,
			B: model.NewLocation(50.0, 20.0),
			C: model.NewRegion([][]float64{
				{-150.0, -150.0},
				{-150.0, 150.0},
				{150.0, 150.0},
				{150.0, -150.0},
				{-150.0, -150.0},
			}),
		},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SpatialModelTest").Where(query.Equal{Column: "SpatialModelTest.B", Value: model.NewLocation(-60.0, 0.0)}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong")
	} else if !res[0].B.IsEqual(model.NewLocation(-60.0, 0.0)) || res[0].A != 1 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}
func TestSpatialModelUpdateCustom(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(SpatialModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("SpatialModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SpatialModelTest{
		{
			A: 1,
			B: model.NewLocation(-60.0, 0.0),
			C: model.NewRegion([][]float64{
				{-100.0, -100.0},
				{-100.0, 100.0},
				{100.0, 100.0},
				{100.0, -100.0},
				{-100.0, -100.0},
			}),
		},
		{
			A: 2,
			B: model.NewLocation(30.0, 10.0),
			C: model.NewRegion([][]float64{
				{-100.0, -100.0},
				{-100.0, 100.0},
				{100.0, 100.0},
				{100.0, -100.0},
				{-100.0, -100.0},
			}),
		},
		{
			A: 3,
			B: model.NewLocation(50.0, 20.0),
			C: model.NewRegion([][]float64{
				{-150.0, -150.0},
				{-150.0, 150.0},
				{150.0, 150.0},
				{150.0, -150.0},
				{-150.0, -150.0},
			}),
		},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	updatedVal := vals[0]
	updatedVal.A = 7
	_, err = db.Model("SpatialModelTest").Where(query.Equal{Column: "B", Value: model.NewLocation(-60.0, 0.0)}).Update(updatedVal)
	if err != nil {
		t.Errorf("Failed to update: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SpatialModelTest").Where(query.Equal{Column: "B", Value: model.NewLocation(-60.0, 0.0)}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong (size)")
	} else if !res[0].B.IsEqual(model.NewLocation(-60.0, 0.0)) || res[0].A != 7 {
		t.Errorf("Query is wrong (record)")
	}

	_, err = db.Execute("DROP TABLE spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestSpatialModelUpdatePrimary(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(SpatialModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("SpatialModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SpatialModelTest{
		{
			A: 1,
			B: model.NewLocation(-60.0, 0.0),
			C: model.NewRegion([][]float64{
				{-100.0, -100.0},
				{-100.0, 100.0},
				{100.0, 100.0},
				{100.0, -100.0},
			}),
		},
		{
			A: 2,
			B: model.NewLocation(30.0, 10.0),
			C: model.NewRegion([][]float64{
				{-100.0, -100.0},
				{-100.0, 100.0},
				{100.0, 100.0},
				{100.0, -100.0},
			}),
		},
		{
			A: 3,
			B: model.NewLocation(50.0, 20.0),
			C: model.NewRegion([][]float64{
				{-150.0, -150.0},
				{-150.0, 150.0},
				{150.0, 150.0},
				{150.0, -150.0},
			}),
		},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	updatedVal := vals[0]
	updatedVal.B = model.NewLocation(-50.0, 15.0)
	_, err = db.Update(updatedVal)
	if err != nil {
		t.Errorf("Failed to update: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SpatialModelTest").Where(query.Equal{Column: "A", Value: 1}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong (size)")
	} else if !res[0].B.IsEqual(model.NewLocation(-50.0, 15.0)) || res[0].A != 1 {
		t.Errorf("%v", res[0].B)
		t.Errorf("Query is wrong (record)")
	}

	_, err = db.Execute("DROP TABLE spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

type SimpleSpatialModelTest struct {
	A int            `atlas:"primarykey;autoincrement"`
	B model.Location `atlas:"default"`
}

func TestSimpleSpatialModelCoverQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(SimpleSpatialModelTest{B: model.NewLocation(40.0, 40.0)})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("SimpleSpatialModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SimpleSpatialModelTest{
		{
			B: model.NewLocation(-60.0, 0.0),
		},
		{
			B: model.NewLocation(30.0, 10.0),
		},
		{
			B: model.NewLocation(50.0, 20.0),
		},
	}

	_, err = db.Model("SimpleSpatialModelTest").Omit("A").Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	_, err = db.Model("SimpleSpatialModelTest").Omit("SimpleSpatialModelTest.B").Create(SimpleSpatialModelTest{A: 100})
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SimpleSpatialModelTest").CoveredBy(model.NewRectRegion(-50.0, 50.0, -50.0, 50.0)).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 3 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE simple_spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

type OrderModelTest struct {
	A int
	B int
	C int
}

func TestModelOrderQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(OrderModelTest{})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("OrderModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []OrderModelTest{
		{A: 1, B: 1, C: 1},
		{A: 1, B: 2, C: 7},
		{A: 1, B: 2, C: 2},
		{A: 1, B: 2, C: 4},
		{A: 1, B: 3, C: 3},
		{A: 2, B: 1, C: 1},
		{A: 2, B: 2, C: 2},
		{A: 2, B: 2, C: 3},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]OrderModelTest, 0)
	err = db.Model("OrderModelTest").Where(query.Equal{Column: "A", Value: "1"}).Limit(3).OrderByCol("B", true).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 3 {
		t.Errorf("Query is wrong")
	} else if res[0].B != 3 || res[2].B != 2 {
		t.Errorf("Query is wrong")
	}

	err = db.Model("OrderModelTest").Where(query.Equal{Column: "A", Value: "1"}).Limit(3).OrderByCol("OrderModelTest.B", true).OrderByCol("C", false).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 3 {
		t.Errorf("Query is wrong")
	} else if res[0].B != 3 || res[0].C != 3 || res[2].B != 2 || res[2].C != 4 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE order_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}

func TestSimpleSpatialModelOrderQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	err = db.RegisterModel(SimpleSpatialModelTest{B: model.NewLocation(40.0, 40.0)})
	if err != nil {
		t.Errorf("Failed to register model: %s\n", err.Error())
	}
	err = db.CreateTable("SimpleSpatialModelTest", true)
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SimpleSpatialModelTest{
		{
			B: model.NewLocation(-60.0, 0.0),
		},
		{
			B: model.NewLocation(30.0, 10.0),
		},
		{
			B: model.NewLocation(50.0, 20.0),
		},
		{
			B: model.NewLocation(40.0, 40.0),
		},
	}

	_, err = db.Model("SimpleSpatialModelTest").Omit("A").Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SimpleSpatialModelTest").OrderByNearestTo(model.NewRectRegion(-10.0, 10.0, -10.0, 10.0), false).Limit(2).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	} else if !res[0].B.IsEqual(vals[1].B) || !res[1].B.IsEqual(vals[2].B) {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE simple_spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}
