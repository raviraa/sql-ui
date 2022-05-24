package qrunner

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"sort"
	"strings"
	"sync"

	"github.com/xo/dburl"
	"github.com/xo/usql/drivers"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/metacmd"
	"github.com/xo/usql/rline"
	"github.com/xo/usql/stmt"
)

type Qrunner struct {
	h   *handler.Handler
	buf *bytes.Buffer
	mu  *sync.Mutex
}

type Result struct {
	Rows    [][]string
	Header  []string
	ResJson string
}

var QrunnerNotInitialized = errors.New("not connected to database")

func New(dsn string) (*Qrunner, error) {
	env.Pset("format", "csv")
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	l, err := rline.New(false, "", "/dev/null")

	if err != nil {
		return nil, err
	}
	h := handler.New(l, usr, wd, true)
	ctx := context.Background()
	if err := h.Open(ctx, dsn); err != nil {
		return nil, err
	}
	q := &Qrunner{
		h:   h,
		buf: &bytes.Buffer{},
		mu:  &sync.Mutex{},
	}
	return q, nil
}

func (q *Qrunner) Dsn() string {
	if q == nil {
		return ""
	}
	return q.h.URL().DSN
}

func (q *Qrunner) Query(ctx context.Context, sqlstr string, outjson bool) (*Result, error) {
	if q == nil {
		return nil, QrunnerNotInitialized
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if outjson {
		env.Pset("format", "json")
		defer env.Pset("format", "csv")
	}
	q.buf.Reset()
	prefix := stmt.FindPrefix(sqlstr, true, true, true)
	log.Println("sql", prefix, sqlstr)

	err := q.h.Execute(ctx, q.buf, metacmd.Option{}, prefix, sqlstr, false)
	if err != nil {
		return nil, err
	}
	if outjson {
		return &Result{ResJson: q.buf.String()}, nil
	}
	return parseCsv(q.buf)
}

type Metatype string

const (
	DescribeTable Metatype = "DescribeTable"
	ListTables    Metatype = "ListTables"
	ListDatabases Metatype = "ListDatabases"
)

func (q *Qrunner) Metacmd(ctx context.Context, cmd Metatype, param string, outjson bool) (*Result, error) {
	if q == nil {
		return nil, QrunnerNotInitialized
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if outjson {
		env.Pset("format", "json")
		defer env.Pset("format", "csv")
	}
	q.buf.Reset()
	m, err := drivers.NewMetadataWriter(ctx, q.h.URL(), q.h.DB(), q.buf)
	if err != nil {
		return nil, err
	}
	switch cmd {
	case DescribeTable:
		err = m.DescribeTableDetails(q.h.URL(), param, false, false)
	case ListTables:
		err = m.ListTables(q.h.URL(), "tvmsE", param, false, false)
	case ListDatabases:
		err = m.ListAllDbs(q.h.URL(), "", false)
	default:
		err = errors.New("unknown meta command: " + string(cmd))
	}
	if err != nil {
		return nil, err
	}
	if outjson {
		return &Result{ResJson: q.buf.String()}, nil
	}
	return parseCsv(q.buf)
}

func (q *Qrunner) QsearchMakeQuery(ctx context.Context, tblname, field, search string) (string, error) {
	if q == nil {
		return "", QrunnerNotInitialized
	}
	tbl, err := q.findTable(ctx, tblname)
	if err != nil {
		return "", err
	}
	return tbl.makeQuery(field, search), nil
}

func (q *Qrunner) Close() error {
	if q == nil {
		return QrunnerNotInitialized
	}
	return q.h.Close()
}

func parseCsv(r io.Reader) (*Result, error) {
	cr := csv.NewReader(r)
	resrows, err := cr.ReadAll()
	var res Result
	if len(resrows) > 0 {
		res.Header = resrows[0]
		res.Rows = resrows[1:]
	}
	return &res, err
}

func Drivers() (*Result, error) {
  buf := &bytes.Buffer{}
	available := drivers.Available()
	fmt.Fprintln(buf, "Name,DSN Alias")
	drvmap := make(map[string]string) // only get unique
	drvkeys := make([]string, 0)
	for drv := range available {
		_, aliases := dburl.SchemeDriverAndAliases(drv)
		drvmap[drv] = strings.Join(aliases, " ") + drvmap[drv]
		drvkeys = append(drvkeys, drv)
	}
	sort.Strings(drvkeys)
	for _, drv := range drvkeys {
		fmt.Fprintf(buf, "%s,%s\n", drv, drvmap[drv])
	}
  return parseCsv(buf)
}
