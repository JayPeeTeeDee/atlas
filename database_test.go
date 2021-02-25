package atlas

import (
	"testing"

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

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO modeltest VALUES (1, 1), (2, 2), (1, 2)")
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

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}

func TestModelSelectQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO modeltest VALUES (1, 1), (2, 2), (1, 2)")
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

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}

func TestModelAllQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO modeltest VALUES (1, 1), (2, 2), (1, 2)")
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

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}

func TestModelSingleInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	val := ModelTest{A: 1, B: 1}

	err = db.Model("ModelTest").Create(val)
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

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}
func TestModelMultiInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []ModelTest{
		{A: 1, B: 1},
		{A: 2, B: 2},
		{A: 1, B: 2},
	}

	err = db.Model("ModelTest").Create(vals)
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

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}

func TestModelSelectInsert(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int DEFAULT 1);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	vals := []ModelTest{
		{A: 1},
		{A: 2},
		{A: 1},
	}

	err = db.Model("ModelTest").Select("A").Create(vals)
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
		query.Equal{Column: "A", Value: "2"},
		query.Equal{Column: "B", Value: "1"},
	}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 1 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}
