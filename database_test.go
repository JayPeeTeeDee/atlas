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
	B int
}

func TestModelSingleQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int);")
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

	res = ModelTest{}
	err = db.Model("ModelTest").Where(query.And{
		query.Equal{Column: "A", Value: "1"},
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int);")
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
	err = db.Model("ModelTest").Select("A").Where(query.Equal{Column: "A", Value: "1"}).All(&other)
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int);")
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int);")
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int);")
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int DEFAULT 1);")
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

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS model_test (a int, b int DEFAULT 1);")
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

	db.RegisterModel(&SpatialModelTest{})
	_, err = db.Execute("CREATE TABLE IF NOT EXISTS spatial_model_test (a int, b geography, c geography);")
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
	err = db.Model("SpatialModelTest").Where(query.Equal{Column: "B", Value: model.NewLocation(-60.0, 0.0)}).All(&res)
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

	db.RegisterModel(&SpatialModelTest{})
	_, err = db.Execute("CREATE TABLE IF NOT EXISTS spatial_model_test (a int, b geography, c geography);")
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

	db.RegisterModel(&SpatialModelTest{})
	_, err = db.Execute("CREATE TABLE IF NOT EXISTS spatial_model_test (a int, b geography, c geography);")
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
	A int `atlas:"primarykey"`
	B model.Location
}

func TestSimpleSpatialModelCoverQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&SimpleSpatialModelTest{})
	_, err = db.Execute("CREATE TABLE IF NOT EXISTS simple_spatial_model_test (a int PRIMARY KEY, b geography);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []SimpleSpatialModelTest{
		{
			A: 1,
			B: model.NewLocation(-60.0, 0.0),
		},
		{
			A: 2,
			B: model.NewLocation(30.0, 10.0),
		},
		{
			A: 3,
			B: model.NewLocation(50.0, 20.0),
		},
	}

	_, err = db.Create(vals)
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]SpatialModelTest, 0)
	err = db.Model("SimpleSpatialModelTest").CoveredBy(model.NewRectRegion(-50.0, 50.0, -50.0, 50.0)).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE simple_spatial_model_test;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
	db.Disconnect()
}
