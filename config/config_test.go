package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoryEntry(t *testing.T) {
	confLocationStr = "/tmp/sql-ui.toml"
	os.Remove(confLocationStr)
	c, err := GetConfig()
	require.Nil(t, err)
	require.Equal(t, 99, c.HistEntryMax)
	q1 := "select * from t1"
	c.AddHistEntry(q1)
	c.AddHistEntry(q1)
	c.AddHistEntry("select * from t2")
	require.Equal(t, 2, len(c.histEntriesMap))
	t1hist, _ := c.histEntriesMap[q1]
	require.Equal(t, 2, t1hist.RunCount)
	c.HistEntryMax = 9

	// SaveConf and valiadate entries
	require.Nil(t, c.SaveConf())
	c2, err := GetConfig()
	require.Nil(t, err)
	require.Equal(t, 9, c2.HistEntryMax)
	require.Equal(t, 2, len(c2.histEntriesMap))
}

func TestHistEntryMax(t *testing.T) {
	confLocationStr = "/tmp/sql-ui.toml"
	os.Remove(confLocationStr)
	c, err := GetConfig()
  c.HistEntryMax = 9
	require.Nil(t, err)
	for i := 0; i < 14; i++ {
		c.AddHistEntry(fmt.Sprintf("select * from t%v", i))
	}
  // these queries should not be removed
	q9 := "select * from t9"
	c.AddHistEntry(q9)
	c.AddHistEntry(q9)
	q8 := "select * from t8"
	c.AddHistEntry(q8)
  // after save, extra queries will be removed
	require.Nil(t, c.SaveConf())
	c2, err := GetConfig()
	require.Nil(t, err)
	require.Equal(t, 9, len(c2.histEntriesMap))
  _, ok := c2.histEntriesMap[q9]
  require.Equal(t, true, ok)
  _, ok = c2.histEntriesMap[q8]
  require.Equal(t, true, ok)
}

func TestHistEntryDsn(t *testing.T) {
	confLocationStr = "/tmp/sql-ui.toml"
	os.Remove(confLocationStr)
	c, err := GetConfig()
  c.ConnectDSNMax = 9
  require.Nil(t, err)
	for i := 0; i < 14; i++ {
		c.AddDSN(fmt.Sprintf("dsn%v", i))
	}
  require.Equal(t, c.ConnectDSNMax, len(c.ConnectDSN))
  // recently added is dsn13
  require.Equal(t, "dsn13", c.ConnectDSN[c.ConnectDSNMax-1])
  // adding existing will move it as last/recent
  c.AddDSN("dsn9")
  fmt.Println(c.ConnectDSN)
  require.Equal(t, "dsn9", c.ConnectDSN[c.ConnectDSNMax-1])

  // save and check
	require.Nil(t, c.SaveConf())
	c2, err := GetConfig()
	require.Nil(t, err)
	require.Equal(t, 9, len(c2.ConnectDSN))
  require.Equal(t, "dsn9", c2.ConnectDSN[c2.ConnectDSNMax-1])
  require.Equal(t, "dsn5", c2.ConnectDSN[0])
  require.Equal(t, "dsn9", c2.DSNsRecent()[0])
  fmt.Println(c2.DSNsRecent())
}
