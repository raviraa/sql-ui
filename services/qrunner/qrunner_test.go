package qrunner

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/xo/usql/drivers/sqlite3"
	"github.com/xo/usql/env"
	// _ "github.com/xo/usql/internal"
)

func TestConn(t *testing.T) {
	// dsn := "moderncsqlite:///tmp/places2.sqlite"
	q := newDb(t)
	defer q.Close()
  log.Println(env.Pall())

	res, err := q.Metacmd(context.Background(), DescribeTable, "t1", false)
	require.Nil(t, err)
	require.Equal(t, res.Rows[0][0], "name")
	require.Equal(t, res.Rows[0][1], "varchar")
	// require.Contains(t, res, `"Name":"val","Type":"INTEGER",`)

	resrows, err := q.Query(context.Background(), "select * from t1", false)
	require.Nil(t, err)
	require.Nil(t, err)
	fmt.Printf("%+v\n", resrows)
	require.Equal(t, "one", resrows.Rows[0][1])
	require.Equal(t, `nine,nine"|"`, resrows.Rows[3][1])

	_, err = q.Query(context.Background(), "select * from t2", false)
	// _, err = q.Query(context.Background(), "select count(*) as count from t2")
	require.NotNil(t, err)
	require.Equal(t, "sqlite3: 1: no such table: t2", err.Error())
}

func newDb(t *testing.T) *Qrunner {
	dsn := "sqlite://testdata/data.sqlite"
	q, err := New(dsn)
	require.Nil(t, err)
	return q
}

func TestDrivers(t *testing.T) {
	q := newDb(t)
	defer q.Close()

	res, err := Drivers()
	require.Nil(t, err)
	require.Equal(t, "sqlite3", res.Rows[0][0])
	require.Equal(t, "sq file sqlite", res.Rows[0][1])
}

func TestQsearch(t *testing.T) {
	q := newDb(t)
	defer q.Close()

	res, err := q.QsearchMakeQuery(context.Background(), "t1", "*text*", "on")
	require.Nil(t, err)
	log.Println("test res", res)
	require.Equal(t, " name LIKE '%on%'", res)
}

func TestMakeQuery(t *testing.T) {
	tbl := Table{Fields: []Field{
		{Name: "f1", Type: "varchar"},
		{Name: "f2", Type: "INTEGER"},
		{Name: "f3", Type: "text"},
	}}
	tests := []struct {
		field, search, want string
	}{
		{"f1", "one", " f1 LIKE '%one%'"},
		{"f1", "fin", " f1 LIKE '%fin%'"},
		{"f1", ">= 3", " f1 >= 3"},
		{"f1", "> 3", " f1 > 3"},
		{"f1", "like foo", " f1 like foo"},
		{"f1", "foo bar", " f1 LIKE '%foo%bar%'"},
		{TextFieldTerm, "foo bar", " f1 LIKE '%foo%bar%' OR  f3 LIKE '%foo%bar%'"},
	}
	for _, tt := range tests {
		got := tbl.makeQuery(tt.field, tt.search)
		require.Equal(t, tt.want, got, tt)
	}
}
