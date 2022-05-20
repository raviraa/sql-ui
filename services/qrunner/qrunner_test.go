package qrunner

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	_ "github.com/xo/usql/drivers/sqlite3"
	// _ "github.com/xo/usql/internal"
)

func TestConn(t *testing.T) {

	// dsn := "moderncsqlite:///tmp/places2.sqlite"
	dsn := "sqlite://testdata/data.sqlite"
	q, err := New(dsn)
	require.Nil(t, err)
	defer q.Close()

	res, err := q.Metacmd(context.Background(), DescribeTable, "t1")
	require.Nil(t, err)
	println(res)
	require.Contains(t, res, `name,varchar,`)
	// require.Contains(t, res, `"Name":"val","Type":"INTEGER",`)

  resrows, err := q.Query(context.Background(), "select * from t1")
	require.Nil(t, err)
	require.Nil(t, err)
  fmt.Printf("%+v\n", resrows)
  require.Equal(t, "name", resrows.Rows[0][1])
  require.Equal(t, `nine,nine"|"`, resrows.Rows[4][1])

	_, err = q.Query(context.Background(), "select * from t2")
	// _, err = q.Query(context.Background(), "select count(*) as count from t2")
	require.NotNil(t, err)
	require.Equal(t, "sqlite3: 1: no such table: t2", err.Error())
}
