package routes

import (
	"fmt"
	"sql-ui/controller"
	"sql-ui/services/qrunner"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo/v4"
)

type (
	tables struct {
		controller.Controller
	}

	tablesData struct {
		Errmsg      string
		PageNumPrev int64
		PageNumNext int64
		Result      *qrunner.Result
		// Table name
		Tname string
		// base url of page with just name param
		Url string
	}
)

// Meta runs schema commands on qrunner db
func (c *tables) Meta(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "query-results"
	page.TemplateExtra = []string{"tables-browse-htmx"}
	page.Data = tablesData{}
	metacmd := ctx.Param("metacmd")
	page.Title = "Schema " + metacmd
	if metacmd == "ListTables" {
		page.Name = "tables-list"
	}

	var errmsg string
	qr, err := c.Container.Qrunner.Metacmd(ctx.Request.Context(), qrunner.Metatype(metacmd), ctx.Query("name"), false)
	if err != nil {
		errmsg = err.Error()
	}
	page.Data = queryData{
		Errmsg: errmsg,
		Result: qr,
	}

	c.RenderPage(*ctx, page)
}

func (c *tables) Browse(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "tables-browse"
	page.HtmxBase = "data"
	page.TemplateExtra = []string{"tables-browse-htmx"}
	tname := ctx.Query("name")
	pageNum, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	perPage := c.Container.Config.PagerSize
	page.Title = "Browse Table " + tname
	query := fmt.Sprintf("select * from %s ", tname)
	var errmsg string
	var err error

	if page.HTMX.Request.TriggerName != "query" {
		page.HtmxBase = "htmx"
	}

	if page.HTMX.Request.TriggerName == "query" {
		page.Name = "tables-browse-htmx"
		field := "*text*" //TODO add in ui
		search := ctx.Query("query")
		whereq, err := c.Container.Qrunner.QsearchMakeQuery(ctx.Request.Context(), tname, field, search)
		if err != nil {
			errmsg = err.Error()
		}
		query += " where " + whereq
	}

	query = fmt.Sprintf("%s limit %v offset %v",
		query, perPage, perPage*int(pageNum))

	var qr *qrunner.Result
	if errmsg == "" { // earlier makequery error
		qr, err = c.Container.Qrunner.Query(ctx.Request.Context(), query, false)
		if err != nil {
			errmsg = query + "; " + err.Error()
		}
	}
	tdata := tablesData{
		Errmsg:      errmsg,
		Result:      qr,
		PageNumNext: pageNum + 1,
		Tname:       tname,
		Url:         "/tables/browse?name=" + tname,
	}
	if pageNum > 0 {
		tdata.PageNumPrev = pageNum - 1
	}
	page.Data = tdata
	c.RenderPage(*ctx, page)
}
