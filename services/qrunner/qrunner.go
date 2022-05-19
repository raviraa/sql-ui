package qrunner

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/user"

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
}

type Result struct {
  Rows [][]string  
}

func New(dsn string) (*Qrunner, error) {
	env.Pset("format", "csv")
  fmt.Printf("%+v\n", env.Pall())
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
	q := &Qrunner{h, &bytes.Buffer{}}
	return q, nil
}

func (q *Qrunner) Query(ctx context.Context, sqlstr string) (string, error) {
	q.buf.Reset()
  log.Println("sql", sqlstr)
  
	prefix := stmt.FindPrefix(sqlstr, true, true, true)
	err := q.h.Execute(ctx, q.buf, metacmd.Option{}, prefix, sqlstr, false)
	return q.buf.String(), err
}

type Metatype string

const (
	DescribeTable Metatype = "DescribeTable"
	ListTables    Metatype = "ListTables"
)

func (q *Qrunner) Metacmd(ctx context.Context, cmd Metatype, param string) (string, error) {
	//drivers
	q.buf.Reset()
	m, err := drivers.NewMetadataWriter(ctx, q.h.URL(), q.h.DB(), q.buf)
	if err != nil {
		return "", err
	}
	switch cmd {
	case DescribeTable:
		err = m.DescribeTableDetails(q.h.URL(), param, false, false)
	case ListTables:
		err = m.ListTables(q.h.URL(), "tvmsE", param, false, false)
	}
	return q.buf.String(), err
}

func (q *Qrunner) Close() error {
	return q.h.Close()
}
