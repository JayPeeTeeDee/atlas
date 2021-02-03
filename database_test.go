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

func TestModelQuery(t *testing.T) {
	db, err := ConnectWithDSN(DBType_Postgres, "postgresql://johnphua:johnphua@localhost/project")
	if err != nil {
		t.Errorf("Failed to connect to db")
	}

	db.RegisterModel(&ModelTest{})

	_, err = db.Execute("CREATE TABLE IF NOT EXISTS modeltest (a int, b int);")
	if err != nil {
		t.Errorf("Failed to create table: %s\n", err.Error())
	}

	_, err = db.Query("INSERT INTO modeltest VALUES (1, 1), (2, 2), (1, 1)")
	if err != nil {
		t.Errorf("Failed to insert rows: %s\n", err.Error())
	}

	res := make([]ModelTest, 0)
	err = db.Model("ModelTest").Where(query.Equal{Column: "a", Value: "1"}).All(&res)
	if err != nil {
		t.Errorf("Failed to query: %s\n", err.Error())
	}

	if len(res) != 2 {
		t.Errorf("Query is wrong")
	}

	_, err = db.Execute("DROP TABLE modeltest;")
	if err != nil {
		t.Errorf("Failed to drop table: %s\n", err.Error())
	}
}
