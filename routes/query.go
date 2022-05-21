package routes

import (
	"sql-ui/controller"
	"sql-ui/services/qrunner"

	"github.com/gin-gonic/gin"
	// "github.com/labstack/echo/v4"
)

type (
	query struct {
		controller.Controller
	}

	queryData struct {
		Errmsg string
		Result *qrunner.Result
	}
)

func (c *query) Get(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "main"
	page.Name = "empty"
	page.Title = "Query"
	page.Data = queryData{}

	c.RenderPage(*ctx, page)
}

func (c *query) Post(ctx *gin.Context) {
	page := controller.NewPage(*ctx, c.Container)
	page.Layout = "htmx"
	page.Name = "query-results"
	page.Title = "Query Results"
	var errmsg string
	c.Container.Query = ctx.PostForm("query")

	qr, err := c.Container.Qrunner.Query(ctx.Request.Context(), c.Container.Query)
	if err != nil {
		errmsg = err.Error()
	}

	page.Data = queryData{
		Errmsg: errmsg,
		Result: qr,
	}

	c.RenderPage(*ctx, page)
}
