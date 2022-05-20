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
		Query  string
		Errmsg string
		Result *qrunner.Result
	}
)

func (c *query) Get(ctx *gin.Context) {
	page := controller.NewPage(*ctx)
	page.Layout = "main"
	page.Name = "empty"
	page.Title = "Query"
	page.Data = queryData{
		Query: "select * from t1",
	}

	c.RenderPage(*ctx, page)
}

func (c *query) Post(ctx *gin.Context) {
	page := controller.NewPage(*ctx)
	page.Layout = "htmx"
	page.Name = "query-results"
	page.Title = "Query Results"
	var errmsg string

	qr, err := c.Container.Qrunner.Query(ctx.Request.Context(), ctx.PostForm("query"))
	if err != nil {
		errmsg = err.Error()
	}

	page.Data = queryData{
		Query:  ctx.PostForm("query"),
		Errmsg: errmsg,
		Result: qr,
	}

	c.RenderPage(*ctx, page)
}
