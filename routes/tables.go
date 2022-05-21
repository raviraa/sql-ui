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
		Errmsg string
		Result *qrunner.Result
	}
)

func (c *tables) List(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "tables-list"
	page.Title = "Tables in Schema"
	page.Data = tablesData{}

  var errmsg string 
  qr, err := c.Container.Qrunner.Metacmd(ctx.Request.Context(), qrunner.ListTables, "")
	if err != nil {
		errmsg = err.Error()
	}
	page.Data = queryData{
		Errmsg: errmsg,
		Result: qr,
	}

	c.RenderPage(*ctx, page)
}

func (c *tables) Describe(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "query-results"
	page.Title = "Schema of Table"
	page.Data = tablesData{}

  var errmsg string 
  qr, err := c.Container.Qrunner.Metacmd(ctx.Request.Context(), qrunner.DescribeTable, ctx.Query("name"))
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
	page.Name = "query-results"
  tname := ctx.Query("name")
  pageNum, _ := strconv.ParseInt(ctx.Query("page"),  10, 64)
  perPage := c.Container.Config.PagerSize
	page.Title = "Browse Table " + tname
  query := fmt.Sprintf("select * from %s limit %v offset %v",
    tname, perPage, perPage*int(pageNum))

  var errmsg string 
	qr, err := c.Container.Qrunner.Query(ctx.Request.Context(),query)
	if err != nil {
		errmsg = err.Error()
	}

	page.Data = tablesData{
    Errmsg: errmsg,
    Result: qr,
  }
	c.RenderPage(*ctx, page)
}
