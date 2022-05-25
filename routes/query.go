package routes

import (
	"github.com/raviraa/sql-ui/controller"
	"github.com/raviraa/sql-ui/services/qrunner"

	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo/v4"
)

type (
	query struct {
		controller.Controller
	}

	queryData struct {
		Errmsg string
		Timing string
		Result *qrunner.Result
	}
)

func (c *query) Get(ctx *gin.Context) {
	page := controller.NewPage(ctx, c.Container)
	page.Layout = "main"
	page.Name = "empty"
	page.Title = "Query"
	page.Data = queryData{}

	c.RenderPage(ctx, page)
}

func (c *query) Post(ctx *gin.Context) {
	page := controller.NewPage(ctx, c.Container)
	page.Layout = "htmx"
	page.Name = "query-results"
	page.Title = "Query Results"
	c.Container.Query = ctx.PostForm("query")
	data := queryData{}

	qr, err := c.Container.Qrunner.Query(ctx.Request.Context(), c.Container.Query, false)
	if err != nil {
		data.Errmsg = err.Error()
	} else {
		data.Timing = qr.Timing
		c.Container.Config.AddHistEntry(c.Container.Query)
	}

	data.Result = qr
	page.Data = data
	c.RenderPage(ctx, page)
}
