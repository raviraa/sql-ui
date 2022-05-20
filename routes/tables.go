package routes

import (
	"sql-ui/controller"
	"sql-ui/services/qrunner"

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

func (c *tables) Get(ctx *gin.Context) {
	page := controller.NewPage(*ctx)
	page.Layout = "main"
	page.Name = "tables"
	page.Title = "tables"
	page.Data = tablesData{
		Errmsg: "testing",
	}

	c.RenderPage(*ctx, page)
}

func (c *tables) Post(ctx *gin.Context) {
}
