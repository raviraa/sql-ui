package qrunner

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Table struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name     string
	Type     string
	Nullable string
	Default  string
}

func (q *Qrunner) findTable(ctx context.Context, tblname string) (*Table, error) {
	jsonres, err := q.Metacmd(ctx, DescribeTable, tblname, true)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	fields := make([]Field, 0)
	if err = json.Unmarshal([]byte(jsonres.ResJson), &fields); err != nil {
		return nil, err
	}
	tbl := &Table{Name: tblname, Fields: fields}
	return tbl, nil
}

var (
	sqlOperators = regexp.MustCompile(`(=|>|<|>=|<=|<>|BETWEEN|LIKE|IN|in|like|between)\s`)
	spaceRe      = regexp.MustCompile(`\s+`)
	// TextFieldTerm is field name from ui, indicating all text fields
	TextFieldTerm = "*text*"
	TextSqlNames  = []string{
		"varchar",
		"text",
		"char",
		"nchar",
		"nvarchar",
		"ntext",
	}
)

// makeQuery is used to make sql query string from field name and search.
// search string words is used to generate like query, unless string contains any sqlOperators
// field can be name of the field, or special type like *text* which matches all fields of TextSqlNames
func (t Table) makeQuery(field, search string) string {
	qfields := []string{}
	if field == TextFieldTerm {
		for _, tblfield := range t.Fields {
			lcname := strings.ToLower(tblfield.Type)
			for _, txtsqlname := range TextSqlNames {
				if strings.Contains(lcname, txtsqlname) {
					qfields = append(qfields, tblfield.Name)
          break
				}
			}
			// if _, ok := TextSqlNames[tblfield.Type]; ok {
			// 	qfields = append(qfields, tblfield.Name)
			// }
		}
	} else {
		qfields = append(qfields, field)
	}
	qsearch := search
	if !sqlOperators.MatchString(search) {
		qsearch = strings.Trim(search, " ")
		qsearch = spaceRe.ReplaceAllString(qsearch, "%")
		qsearch = "LIKE '%" + qsearch + "%'"
	}

	query := ""
	for idx, qfield := range qfields {
		if idx != 0 {
			query += " OR "
		}
		query += fmt.Sprintf(" %s %s", qfield, qsearch)
	}
	log.Println(field, qfields, qsearch, query)
	return query
}
